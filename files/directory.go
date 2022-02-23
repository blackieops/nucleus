package files

import (
	"path/filepath"
	"time"

	"go.b8s.dev/nucleus/auth"
)

type Directory struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Name of the directory
	Name string

	// Projected full path to this directory including its parents
	FullName string

	// Optional parent directory to support hierarchy
	ParentID *int
	Parent   *Directory `gorm:"constraint:OnDelete:CASCADE;"`

	// Owner of this directory
	UserID int
	User   auth.User
}

func (d *Directory) SetNames(name string) {
	d.Name = name
	if d.Parent == nil {
		d.FullName = name
		return
	}
	d.FullName = filepath.Join(d.Parent.FullName, name)
}
