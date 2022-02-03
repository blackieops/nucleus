package auth

import (
	"com.blackieops.nucleus/data"
	"time"
)

type User struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

func FindUserByUsername(ctx *data.Context, id string) (*User, error) {
	var user *User
	err := ctx.DB.Where("username = ?", id).First(&user).Error
	return user, err
}

func CreateUser(ctx *data.Context, u *User) (*User, error) {
	err := ctx.DB.Create(&u).Error
	return u, err
}
