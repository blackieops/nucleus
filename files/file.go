package files

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"time"

	"go.b8s.dev/nucleus/auth"
)

// File represents a data object in a storage backend
type File struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Name is the basename of the file
	Name string

	// Parent is an optional association to a Directory
	Parent   *Directory `gorm:"constraint:OnDelete:CASCADE;"`
	ParentID *uint

	// User is the owner of this file
	User   auth.User
	UserID uint

	// FullName is a projection of the full path of all parent directories
	FullName string `gorm:"index"`

	// Size is the file size in bytes
	Size int64

	// Digest is a SHA-1 hash of the file content
	Digest string
}

// SetNames sets both the Name and generates the FullName based on the parent,
// if it is set.
func (f *File) SetNames(name string) {
	f.Name = name
	if f.Parent == nil {
		f.FullName = name
		return
	}
	f.FullName = filepath.Join(f.Parent.FullName, name)
}

// SetDigest will calculate the SHA-1 digest of the given byte array content
// and assign the result to `Digest`.
func (f *File) SetDigest(content []byte) {
	digest := sha1.Sum(content)
	f.Digest = fmt.Sprintf("%x", digest[:])
}
