package files

import (
	"io/fs"

	"com.blackieops.nucleus/auth"
)

type StorageBackend interface {
	List(*auth.User, string) []fs.FileInfo
	Stat(*auth.User, string) (fs.FileInfo, error)
	ReadFile(*File) ([]byte, error)
}
