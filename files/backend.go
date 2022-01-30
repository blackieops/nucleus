package files

import (
	"io/fs"

	"com.blackieops.nucleus/auth"
)

type StorageBackend interface {
	List(*auth.User, string) []fs.FileInfo
	Stat(*auth.User, string) (fs.FileInfo, error)
	ReadFile(*auth.User, *File) ([]byte, error)
	FileDigest(*auth.User, *File) (string, error)
}
