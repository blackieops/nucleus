package nxc

import (
	"com.blackieops.nucleus/data"
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

func CreateNextcloudAuthSession(ctx *data.Context) (*NextcloudAuthSession, error) {
	session := &NextcloudAuthSession{
		PollToken:  generateNextcloudToken(64),
		LoginToken: generateNextcloudToken(64),
	}
	err := ctx.DB.Create(session).Error
	return session, err
}

func FindNextcloudAuthSessionByPollToken(ctx *data.Context, token string) (*NextcloudAuthSession, error) {
	var session *NextcloudAuthSession
	err := ctx.DB.Where("poll_token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return session, nil
}

func FindNextcloudAuthSessionByLoginToken(ctx *data.Context, token string) (*NextcloudAuthSession, error) {
	var session *NextcloudAuthSession
	err := ctx.DB.Where("login_token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return session, nil
}

func DestroyNextcloudAuthSession(ctx *data.Context, session *NextcloudAuthSession) error {
	return ctx.DB.Delete(session).Error
}
