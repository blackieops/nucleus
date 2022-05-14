package files

import (
	"path/filepath"
	"time"

	"go.b8s.dev/nucleus/auth"
)

// Directory represents a "folder" or "prefix" in a storage backend.
type Directory struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Name of the directory
	Name string

	// FullName is the projected path to this directory including its parents
	FullName string

	// Parent is an optional parent directory to support hierarchy
	Parent   *Directory `gorm:"constraint:OnDelete:CASCADE;"`
	ParentID *int

	// User is the owner of this directory
	User   auth.User
	UserID int
}

// SetNames sets both the `Name` and generates the `FullName` based on the
// parent, if it is set.
func (d *Directory) SetNames(name string) {
	d.Name = name
	if d.Parent == nil {
		d.FullName = name
		return
	}
	d.FullName = filepath.Join(d.Parent.FullName, name)
}
