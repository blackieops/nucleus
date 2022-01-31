package files

import (
	"crypto/sha1"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/data"
)

type FilesystemBackend struct {
	DBContext     *data.Context
	StoragePrefix string
}

func (b *FilesystemBackend) List(user *auth.User, path string) []fs.FileInfo {
	entries, err := ioutil.ReadDir(b.userStoragePath(user, path))
	if err != nil {
		panic(err)
	}
	return entries
}

func (b *FilesystemBackend) Stat(user *auth.User, path string) (fs.FileInfo, error) {
	fullPath := b.userStoragePath(user, path)
	return os.Stat(fullPath)
}

func (b *FilesystemBackend) ReadFile(user *auth.User, file *File) ([]byte, error) {
	return ioutil.ReadFile(b.userStoragePath(user, file.FullName))
}

func (b *FilesystemBackend) FileDigest(user *auth.User, file *File) (string, error) {
	fileBytes, err := b.ReadFile(user, file)
	if err != nil {
		return "", err
	}
	digest := sha1.Sum(fileBytes)
	return fmt.Sprintf("%x", digest[:]), nil
}

func (b *FilesystemBackend) storagePath(path *string) string {
	if path == nil {
		return b.StoragePrefix
	} else {
		return b.StoragePrefix + string(os.PathSeparator) + *path
	}
}

func (b *FilesystemBackend) userStoragePath(user *auth.User, path string) string {
	sep := string(os.PathSeparator)
	filesBasePath := user.Username + sep + "files" + sep + path
	return b.storagePath(&filesBasePath)
}
