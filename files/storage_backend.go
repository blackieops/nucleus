package files

import (
	"io/fs"

	"go.b8s.dev/nucleus/auth"
)

// StorageBackend represents the interface required to persist files.
type StorageBackend interface {
	// List will return all entries at the given path
	List(*auth.User, string) ([]fs.FileInfo, error)

	// Stat will return the FileInfo for the given entry
	Stat(*auth.User, string) (fs.FileInfo, error)

	// ReadFile will return the file content as bytes.
	ReadFile(*auth.User, *File) ([]byte, error)

	// WriteFile will write the given content for the given File.
	WriteFile(*auth.User, *File, []byte) error

	// CreateDirectory will create the directory in the backend.
	CreateDirectory(*auth.User, *Directory) error

	// DeletePath will recursively delete the entries under the given path
	DeletePath(*auth.User, string) error

	// RenamePath will rename the entry at the given path, and if it has any,
	// also its children.
	RenamePath(*auth.User, string, string) error

	CreateChunkDirectory(*auth.User, string) error
	WriteChunk(*auth.User, string, []byte) error
	ReconstructChunks(*auth.User, string, string) error
	DeleteChunkDirectory(*auth.User, string) error
}
