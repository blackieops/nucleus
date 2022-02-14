package files

import (
	"os"
	"time"

	"com.blackieops.nucleus/auth"
)

type File struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Basename of the file
	Name string

	// Optional association to a Directory
	ParentID *uint
	Parent   *Directory `gorm:"constraint:OnDelete:CASCADE;"`

	// Association to the user who owns this file
	UserID uint
	User   auth.User

	// Projected full path of all parent directories for easier lookups
	FullName string `gorm:"index"`

	// File size in bytes
	Size int64

	// SHA-1 hash of the file content for use as an etag or similar cache key
	Digest string
}

func (f *File) SetNames(name string) {
	f.Name = name
	if f.Parent == nil {
		f.FullName = name
		return
	}
	f.FullName = f.Parent.FullName + string(os.PathSeparator) + name
}
