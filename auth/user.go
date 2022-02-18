package auth

import (
	"crypto/md5"
	"fmt"
	"time"

	"com.blackieops.nucleus/data"
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

func (u *User) AvatarURL(size int) string {
	return fmt.Sprintf(
		"https://www.gravatar.com/avatar/%x?s=%s",
		md5.Sum([]byte(u.EmailAddress)),
		fmt.Sprint(size),
	)
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

func UpdateUser(ctx *data.Context, u *User) (*User, error) {
	err := ctx.DB.Save(u).Error
	return u, err
}

func DeleteUser(ctx *data.Context, u *User) error {
	return ctx.DB.Delete(u).Error
}
