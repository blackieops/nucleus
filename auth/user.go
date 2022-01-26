package auth

import (
	"com.blackieops.nucleus/data"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"index"`
	Name         string
	EmailAddress string `gorm:"index"`
}

func FindAllUsers(ctx *data.Context) []*User {
	var users []*User
	ctx.DB.Find(&users, User{})
	return users
}

func FindUser(ctx *data.Context, id int) *User {
	var user *User
	ctx.DB.First(&user, id)
	return user
}
