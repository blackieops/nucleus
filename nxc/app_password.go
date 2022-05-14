package nxc

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"

	"go.b8s.dev/nucleus/auth"
	"go.b8s.dev/nucleus/data"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// NextcloudAppPassword represents a connection with a Nextcloud client. Sort of
// an "oauth token lite™".
type NextcloudAppPassword struct {
	gorm.Model
	PasswordDigest string `gorm:"index"`
	UserID         int
	User           auth.User
}

func ListNextcloudAppPasswordsForUser(
	ctx *data.Context,
	user *auth.User,
) ([]*NextcloudAppPassword, error) {
	var passwords []*NextcloudAppPassword
	err := ctx.DB.Where("user_id = ?", user.ID).Find(&passwords).Error
	return passwords, err
}

func FindNextcloudAppPassword(
	ctx *data.Context,
	user *auth.User,
	id uint,
) (*NextcloudAppPassword, error) {
	var password *NextcloudAppPassword
	err := ctx.DB.Where("user_id = ? and id = ?", user.ID, id).First(&password).Error
	return password, err
}

func CreateNextcloudAppPassword(
	ctx *data.Context,
	session *NextcloudAuthSession,
	user *auth.User,
) (*NextcloudAppPassword, error) {
	password := generateNextcloudToken(32)
	// We use a low cost here because these are very long random strings, and
	// these are easy and safe to rotate in the face of a disasterous security
	// event. It's mostly just a deterrent to buy time to rotate everything,
	// not to protect an entire account forever.
	digest, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		panic(err)
	}
	appPassword := &NextcloudAppPassword{
		PasswordDigest: string(digest),
		User:           *user,
	}

	ctx.DB.Create(appPassword)

	// We need to prefix the ID to the password so we have a way to look up
	// which AppPassword this is when we compare the hash later.
	session.RawAppPassword = fmt.Sprint(appPassword.ID) + "-" + password
	session.Username = user.Username

	ctx.DB.Save(session)

	return appPassword, nil
}

func FindNextcloudAppPasswordByPassword(c *data.Context, composite string) (*NextcloudAppPassword, error) {
	bits := strings.Split(composite, "-")
	var appPassword *NextcloudAppPassword
	err := c.DB.Preload("User").Where("id = ?", bits[0]).First(&appPassword).Error
	if err != nil {
		return &NextcloudAppPassword{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(appPassword.PasswordDigest), []byte(bits[1]))
	if err != nil {
		return &NextcloudAppPassword{}, err
	}
	return appPassword, err
}

func DeleteNextcloudAppPassword(c *data.Context, p *NextcloudAppPassword) error {
	return c.DB.Delete(p).Error
}

func generateNextcloudToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
