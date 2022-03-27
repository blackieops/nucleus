package files

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.b8s.dev/nucleus/auth"
)

type FilesystemBackend struct {
	StoragePrefix string
}

func (b *FilesystemBackend) List(user *auth.User, path string) ([]fs.FileInfo, error) {
	entries, err := ioutil.ReadDir(b.userStoragePath(user, path))
	if err != nil {
		return []fs.FileInfo{}, err
	}
	return entries, nil
}

func (b *FilesystemBackend) Stat(user *auth.User, path string) (fs.FileInfo, error) {
	fullPath := b.userStoragePath(user, path)
	return os.Stat(fullPath)
}

func (b *FilesystemBackend) ReadFile(user *auth.User, file *File) ([]byte, error) {
	return ioutil.ReadFile(b.userStoragePath(user, file.FullName))
}

func (b *FilesystemBackend) WriteFile(user *auth.User, file *File, contents []byte) error {
	return ioutil.WriteFile(b.userStoragePath(user, file.FullName), contents, 0644)
}

func (b *FilesystemBackend) CreateDirectory(user *auth.User, dir *Directory) error {
	err := os.MkdirAll(b.userStoragePath(user, dir.FullName), 0755)
	if err != nil {
		return err
	}
	return nil
}

func (b *FilesystemBackend) DeletePath(user *auth.User, path string) error {
	if path == "" || path == "." {
		return errors.New("Refusing to delete all files from storage. FullName is empty.")
	}
	return os.RemoveAll(b.userStoragePath(user, path))
}

func (b *FilesystemBackend) RenamePath(user *auth.User, src string, dest string) error {
	return os.Rename(b.userStoragePath(user, src), b.userStoragePath(user, dest))
}

func (b *FilesystemBackend) CreateChunkDirectory(user *auth.User, name string) error {
	err := os.Mkdir(b.userUploadsPath(user, name), 0755)
	if err != nil {
		return err
	}
	return nil
}

func (b *FilesystemBackend) WriteChunk(user *auth.User, name string, contents []byte) error {
	return ioutil.WriteFile(b.userUploadsPath(user, name), contents, 0644)
}

func (b *FilesystemBackend) ReconstructChunks(user *auth.User, srcDir string, destPath string) error {
	destFile, err := os.Create(b.userStoragePath(user, destPath))
	if err != nil {
		return err
	}
	defer destFile.Close()
	uploadDir := b.userUploadsPath(user, srcDir)
	chunks, err := ioutil.ReadDir(uploadDir)
	if err != nil {
		return err
	}
	for _, chunk := range chunks {
		chunkBytes, err := ioutil.ReadFile(uploadDir + "/" + chunk.Name())
		if err != nil {
			return err
		}
		_, err = destFile.Write(chunkBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *FilesystemBackend) DeleteChunkDirectory(user *auth.User, dirName string) error {
	return os.RemoveAll(b.userUploadsPath(user, dirName))
}

func (b *FilesystemBackend) storagePath(path *string) string {
	if path == nil {
		return b.StoragePrefix
	} else {
		return filepath.Join(b.StoragePrefix, *path)
	}
}

func (b *FilesystemBackend) userStoragePath(user *auth.User, path string) string {
	filesBasePath := filepath.Join(user.Username, "files", path)
	return b.storagePath(&filesBasePath)
}

func (b *FilesystemBackend) userUploadsPath(user *auth.User, path string) string {
	uploadsBasePath := filepath.Join(user.Username, "uploads", path)
	return b.storagePath(&uploadsBasePath)
}
