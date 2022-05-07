package files

import (
	"io/fs"
	"time"
	"strings"

	"go.b8s.dev/nucleus/auth"
)


// MemoryFile conforms to the fs.FileInfo interface.
type MemoryFile struct {
	name string
	size int64
	isDir bool
	modTime time.Time
}

func (f MemoryFile) Name() string {
	return f.name
}
func (f MemoryFile) Size() int64 {
	return f.size
}
func (f MemoryFile) IsDir() bool {
	return f.isDir
}

func (f MemoryFile) ModTime() time.Time {
	return f.modTime
}

func (f MemoryFile) Mode() fs.FileMode {
	return 0
}

func (f MemoryFile) Sys() interface{} {
	return nil
}

type MemoryStorageBackend struct {
	entries map[string][]byte
}

func (b *MemoryStorageBackend) List(u *auth.User, p string) ([]fs.FileInfo, error) {
	files := []fs.FileInfo{}
	for key := range b.entries {
		if strings.HasPrefix(key, p) {
			files = append(files, b.entryToMemoryFile(key))
		}
	}
	return files, nil
}

func (b *MemoryStorageBackend) Stat(u *auth.User, path string) (fs.FileInfo, error) {
	return b.entryToMemoryFile(path), nil
}

func (b *MemoryStorageBackend) ReadFile(u *auth.User, f *File) ([]byte, error) {
	return b.entries[f.FullName], nil
}

func (b *MemoryStorageBackend) WriteFile(u *auth.User, f *File, content []byte) error {
	b.initEntries()
	b.entries[f.FullName] = content
	return nil
}

func (b *MemoryStorageBackend) DeletePath(u *auth.User, path string) error {
	b.initEntries()
	for key := range b.entries {
		if strings.HasPrefix(key, path) {
			delete(b.entries, key)
		}
	}
	return nil
}

func (b *MemoryStorageBackend) RenamePath(u *auth.User, src, dest string) error {
	b.initEntries()
	copy(b.entries[dest], b.entries[src])
	delete(b.entries, src)
	return nil
}

func (b *MemoryStorageBackend) initEntries() {
	if b.entries == nil {
		b.entries = make(map[string][]byte)
	}
}

func (b *MemoryStorageBackend) entryToMemoryFile(key string) MemoryFile {
	return MemoryFile{
		name: key,
		size: int64(len(b.entries[key])),
		isDir: false,
		modTime: time.Now(),
	}
}
