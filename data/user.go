package data

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username              string `gorm:"index"`
	Name                  string
	EmailAddress          string `gorm:"index"`
	NextcloudAppPasswords []NextcloudAppPassword
}
