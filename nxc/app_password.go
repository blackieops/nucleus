package nxc

import (
	"fmt"
	"encoding/hex"
	"math/rand"
	"strings"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/data"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type NextcloudAppPassword struct {
	gorm.Model
	PasswordDigest string `gorm:"index"`
	UserID         int
	User           auth.User
}

func CreateNextcloudAppPassword(
	ctx *data.Context,
	session *NextcloudAuthSession,
	user *auth.User,
) (*NextcloudAppPassword, error) {
	password := generateNextcloudToken(64)
	// TODO: parameterize cost into config value?
	digest, err := bcrypt.GenerateFromPassword([]byte(password), 10)
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
	session.RawAppPassword = fmt.Sprint(appPassword.ID)+"-"+password
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

func generateNextcloudToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
