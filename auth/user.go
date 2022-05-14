package auth

import (
	"crypto/md5"
	"fmt"
	"time"

	"go.b8s.dev/nucleus/data"
)

// User represents a person using the system. A user can be authenticated via
// its `Credential`s to gain access.
type User struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Username     string `gorm:"uniqueIndex"`
	Name         string
	EmailAddress string `gorm:"uniqueIndex"`
	Credentials  []*Credential
}

// AvatarURL returns the URL to the user's avatar image of the given size.
func (u *User) AvatarURL(size int) string {
	return fmt.Sprintf(
		"https://www.gravatar.com/avatar/%x?s=%s",
		md5.Sum([]byte(u.EmailAddress)),
		fmt.Sprint(size),
	)
}

// FindAllUsers returns a list of all users in the system.
func FindAllUsers(ctx *data.Context) []*User {
	var users []*User
	ctx.DB.Find(&users, User{})
	return users
}

// FindUser looks up a user by its internal ID.
func FindUser(ctx *data.Context, id uint) (*User, error) {
	var user *User
	err := ctx.DB.First(&user, id).Error
	return user, err
}

// FindUserByUsername looks up a user by its username.
func FindUserByUsername(ctx *data.Context, id string) (*User, error) {
	var user *User
	err := ctx.DB.Where("username = ?", id).First(&user).Error
	return user, err
}

// CreateUser persists the given user. The returned user will have
// auto-generated fields such as timestamps and IDs populated.
func CreateUser(ctx *data.Context, u *User) (*User, error) {
	err := ctx.DB.Create(&u).Error
	return u, err
}

// UpdateUser saves the changes to the given user.
func UpdateUser(ctx *data.Context, u *User) (*User, error) {
	err := ctx.DB.Save(u).Error
	return u, err
}

// DeleteUser deletes the user immediately.
func DeleteUser(ctx *data.Context, u *User) error {
	return ctx.DB.Delete(u).Error
}
