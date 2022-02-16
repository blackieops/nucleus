package auth

import (
	"errors"
	"time"

	"com.blackieops.nucleus/data"
	"golang.org/x/crypto/bcrypt"
)

const CredentialTypePassword int = 0

type Credential struct {
	ID         uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Type       int    `gorm:"index,not null,default:0"`
	Data       string `gorm:"-"`
	DataDigest string `gorm:"not null"`
	UserID     uint
	User       User
}

// Find all the Credentials for the given User.
func FindUserCredentials(c *data.Context, user *User) ([]*Credential, error) {
	var credentials []*Credential
	err := c.DB.Where("user_id = ?", user.ID).Find(&credentials).Error
	return credentials, err
}

// Create a new Credential for the given user. Takes a `type`, which is an enum
// value (see the consts under `auth` for available types); and `data` as a
// string to be hashed (eg., a password).
func CreateCredential(c *data.Context, user *User, credential *Credential) (*Credential, error) {
	if credential.Data == "" {
		return nil, errors.New("Data cannot be empty.")
	}
	if credential.DataDigest != "" {
		return nil, errors.New("DataDigest cannot be set directly.")
	}

	// TODO: parameterize cost in config?
	dataDigest, err := bcrypt.GenerateFromPassword([]byte(credential.Data), 12)
	credential.DataDigest = string(dataDigest)
	credential.User = *user
	credential.Type = CredentialTypePassword

	err = c.DB.Create(credential).Error
	return credential, err
}

func UpdateCredential(ctx *data.Context, c *Credential, data string) (*Credential, error) {
	// TODO: parameterize cost in config?
	dataDigest, err := bcrypt.GenerateFromPassword([]byte(data), 12)
	if err != nil {
		return c, err
	}
	c.DataDigest = string(dataDigest)
	err = ctx.DB.Save(c).Error
	return c, err
}

// Deletes a credential from the database immediately.
func DeleteCredential(c *data.Context, credential *Credential) error {
	return c.DB.Where("id = ?", credential.ID).Delete(credential).Error
}

// Checks if the given string matches the hashed data on the Credential.
// Primarily useful for passwords. Returns error if they do not match.
func ValidateCredential(c *Credential, data string) error {
	return bcrypt.CompareHashAndPassword([]byte(c.DataDigest), []byte(data))
}

// RIP generics
// Filter a slice of Credentials to find the first of a particular type.
func FilterFirstCredentialOfType(cs []*Credential, t int) (*Credential, error) {
	for _, c := range cs {
		if c.Type == t {
			return c, nil
		}
	}
	return nil, errors.New("No credential found for type.")
}
