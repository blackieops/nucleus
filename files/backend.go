package files

import (
	"io"
	"io/fs"

	"com.blackieops.nucleus/auth"
)

type StorageBackend interface {
	List(*auth.User, string) []fs.FileInfo
	Stat(*auth.User, string) (fs.FileInfo, error)
	ReadFile(*auth.User, *File) ([]byte, error)
	ReaderFile(*auth.User, *File) (io.ReadCloser, error)
	WriteFile(*auth.User, *File, []byte) error
	CreateDirectory(*auth.User, *Directory) error
	DeletePath(*auth.User, string) error
	RenamePath(*auth.User, string, string) error
	CreateChunkDirectory(*auth.User, string) error
	WriteChunk(*auth.User, string, []byte) error
	ReconstructChunks(*auth.User, string, string) error
	DeleteChunkDirectory(*auth.User, string) error
}
