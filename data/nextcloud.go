package data

import (
	"encoding/hex"
	"errors"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type NextcloudAuthSession struct {
	gorm.Model
	PollToken  string `gorm:"index"`
	LoginToken string `gorm:"index"`

	// An unfortunate compromise -- we need to serialize the app password in
	// plain text to the client when it polls. Since the actions are in
	// entirely separate sessions (the user granting, and then the client
	// polling), we need to temporarily store the raw app password somewhere.
	//
	// Since this record gets destroyed immediately after serialising, this is
	// "probably fine". And since industry practice for auth tokens is to just
	// yolo store them plain text, we're probably still ahead of the game. But
	// still, it's a compromise.
	RawAppPassword string

	// Since we just need to serialize the username for the Poll endpoints, we
	// just cache it here to avoid having to do two queries to get to the User.
	Username string
}

type NextcloudAppPassword struct {
	gorm.Model
	PasswordDigest string `gorm:"index"`
	UserID         int
	User           User
}

func CreateNextcloudAuthSession(ctx *Context) *NextcloudAuthSession {
	session := &NextcloudAuthSession{
		PollToken:  generateNextcloudToken(64),
		LoginToken: generateNextcloudToken(64),
	}
	ctx.DB.Create(session)
	return session
}

func FindNextcloudAuthSessionByPollToken(ctx *Context, token string) (*NextcloudAuthSession, error) {
	var session *NextcloudAuthSession
	ctx.DB.Where("poll_token = ?", token).First(&session)

	if session.ID != 0 {
		return session, nil
	} else {
		return &NextcloudAuthSession{}, errors.New("Could not find session.")
	}
}

func FindNextcloudAuthSessionByLoginToken(ctx *Context, token string) (*NextcloudAuthSession, error) {
	var session *NextcloudAuthSession
	ctx.DB.Where("login_token = ?", token).First(&session)

	if session.ID != 0 {
		return session, nil
	} else {
		return &NextcloudAuthSession{}, errors.New("Could not find session.")
	}
}

func DestroyNextcloudAuthSession(ctx *Context, session *NextcloudAuthSession) error {
	ctx.DB.Delete(session)
	// TODO: errors?
	return nil
}

func CreateNextcloudAppPassword(ctx *Context, session *NextcloudAuthSession, user *User) (*NextcloudAppPassword, error) {
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
