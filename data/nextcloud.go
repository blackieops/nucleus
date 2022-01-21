package data

import (
	"gorm.io/gorm"
)

type NextcloudAuthSession struct {
	gorm.Model
	PollToken     string `gorm:"index"`
	LoginToken    string `gorm:"index"`
	AppPasswordID int
	AppPassword   NextcloudAppPassword
}

type NextcloudAppPassword struct {
	gorm.Model
	PasswordDigest string `gorm:"index"`
	UserID         int
	User           User
}
