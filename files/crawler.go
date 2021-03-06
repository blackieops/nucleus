package files

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"go.b8s.dev/nucleus/auth"
	"go.b8s.dev/nucleus/data"
)

type Crawler struct {
	DBContext *data.Context
	Backend   StorageBackend
}

func (c *Crawler) ReindexAll() {
	users := auth.FindAllUsers(c.DBContext)

	for _, user := range users {
		c.IndexUserFiles(user, nil)
	}
}

func (c *Crawler) IndexUserFiles(user *auth.User, currentDir *Directory) {
	var wg sync.WaitGroup
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
				wg.Add(1)
				go func() {
					defer wg.Done()
					c.IndexUserFiles(user, directory)
				}()
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

	wg.Wait()
}

// Query for the file in the storage backend and return an unpersisted struct
// with some of the basic file metadata filled out.
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
