package files

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"go.b8s.dev/nucleus/auth"
	"go.b8s.dev/nucleus/data"
)

// crawlers is a channel that acts as a counting semaphore to prevent
// goroutines going too wild when indexing a large directory tree.
// TODO: it would be nice if this were configurable.
var crawlers = make(chan bool, 32)

// wg keeps a counter of all goroutines currently in-flight to ensure they all
// finish before the program exits.
var wg sync.WaitGroup

// Crawler provides methods to index the contents of a storage backend.
type Crawler struct {
	DBContext *data.Context
	Backend   StorageBackend
}

// ReindexAll will crawl the entire storage backend for all users in the system
// and index all the files it finds.
func (c *Crawler) ReindexAll() {
	users := auth.FindAllUsers(c.DBContext)
	for _, user := range users {
		c.IndexUserFiles(user, nil)
	}
	wg.Wait()
}

// IndexUserFiles will index all files for the given user in the storage
// backend.
func (c *Crawler) IndexUserFiles(user *auth.User, currentDir *Directory) {
	wg.Add(1)
	defer wg.Done()
	crawlers <- true

	var entries []fs.FileInfo
	var err error
	if currentDir == nil {
		entries, err = c.Backend.List(user, "")
	} else {
		entries, err = c.Backend.List(user, currentDir.FullName)
	}
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		fmt.Println("Indexing file: ", entry.Name())

		if entry.IsDir() {
			newDir := c.DiscoverDir(user, currentDir, entry.Name())
			directory, err := CreateDir(c.DBContext, newDir)
			if err == nil {
				go c.IndexUserFiles(user, directory)
			} else {
				fmt.Println("Failed to discover directory in backend: ", entry.Name())
			}
		} else {
			fileData, err := c.DiscoverFile(user, currentDir, entry.Name())
			if err == nil {
				CreateFile(c.DBContext, fileData)
			} else {
				fmt.Println("Failed to discover file in backend: ", entry.Name())
			}
		}
	}

	<- crawlers
}

// DiscoverFile queries for the file in the storage backend and returns an
// unpersisted struct with some of the basic file metadata filled out.
func (c *Crawler) DiscoverFile(user *auth.User, dir *Directory, name string) (*File, error) {
	var path string
	if dir == nil {
		path = name
	} else {
		path = filepath.Join(dir.FullName, name)
	}
	fileStat, err := c.Backend.Stat(user, path)
	if err != nil {
		return nil, err
	}
	fileEntity := &File{
		Name:     fileStat.Name(),
		FullName: path,
		Size:     fileStat.Size(),
		User:     *user,
	}
	if dir != nil {
		fileEntity.Parent = dir
	}
	fileContents, err := c.Backend.ReadFile(user, fileEntity)
	if err != nil {
		return nil, err
	}
	fileEntity.SetDigest(fileContents)
	return fileEntity, nil
}

// DiscoverDir finds the directory in the storage backend and returns an
// unpersisted struct with some basic metadata populated.
func (c *Crawler) DiscoverDir(user *auth.User, parent *Directory, name string) *Directory {
	directory := &Directory{
		Name: name,
		User: *user,
	}
	if parent != nil {
		directory.FullName = parent.FullName + "/" + name
		directory.Parent = parent
	} else {
		directory.FullName = name
	}
	return directory
}
