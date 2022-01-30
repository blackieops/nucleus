package files

import (
	"com.blackieops.nucleus/auth"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model

	// Basename of the file
	Name string

	// Optional association to a Directory
	ParentID *int
	Parent   *Directory

	// Association to the user who owns this file
	UserID int
	User   auth.User

	// Projected full path of all parent directories for easier lookups
	FullName string `gorm:"index"`

	// File size in bytes
	Size int64

	// SHA-1 hash of the file content for use as an etag or similar cache key
	Digest string
}

type Directory struct {
	gorm.Model

	// Name of the directory
	Name string

	// Projected full path to this directory including its parents
	FullName string

	// Optional parent directory to support hierarchy
	ParentID *int
	Parent   *Directory

	// Owner of this directory
	UserID int
	User   auth.User
}
