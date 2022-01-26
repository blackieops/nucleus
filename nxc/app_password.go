package nxc

import (
	"encoding/hex"
	"math/rand"

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

	session.RawAppPassword = password
	session.Username = user.Username

	ctx.DB.Create(appPassword)
	ctx.DB.Save(session)

	return appPassword, nil
}

func generateNextcloudToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
