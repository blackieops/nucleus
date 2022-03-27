package files

import (
	"io/fs"

	"go.b8s.dev/nucleus/auth"
)

type StorageBackend interface {
	List(*auth.User, string) ([]fs.FileInfo, error)
	Stat(*auth.User, string) (fs.FileInfo, error)
	ReadFile(*auth.User, *File) ([]byte, error)
	WriteFile(*auth.User, *File, []byte) error
	CreateDirectory(*auth.User, *Directory) error
	DeletePath(*auth.User, string) error
	RenamePath(*auth.User, string, string) error
	CreateChunkDirectory(*auth.User, string) error
	WriteChunk(*auth.User, string, []byte) error
	ReconstructChunks(*auth.User, string, string) error
	DeleteChunkDirectory(*auth.User, string) error
}
