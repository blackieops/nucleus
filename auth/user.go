package auth

import (
	"com.blackieops.nucleus/data"
	"time"
)

type User struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Username     string `gorm:"uniqueIndex"`
	Name         string
	EmailAddress string `gorm:"uniqueIndex"`
	Credentials  []*Credential
}

func FindAllUsers(ctx *data.Context) []*User {
	var users []*User
	ctx.DB.Find(&users, User{})
	return users
}

func FindUser(ctx *data.Context, id uint) (*User, error) {
	var user *User
	err := ctx.DB.First(&user, id).Error
	return user, err
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
