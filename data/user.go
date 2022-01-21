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

func FindAllUsers(ctx *Context) []*User {
	var users []*User
	ctx.DB.Find(&users, User{})
	return users
}

func FindUser(ctx *Context, id int) *User {
	var user *User
	ctx.DB.First(&user, id)
	return user
}
