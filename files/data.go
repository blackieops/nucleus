package files

import (
	"path/filepath"

	"go.b8s.dev/nucleus/auth"
	"go.b8s.dev/nucleus/data"
)

// CompositeListing represents a combined Directory and File enumeration.
type CompositeListing struct {
	// Parent is the directory for under which these files and directories sit
	Parent *Directory

	// The files in this directory
	Files []*File

	// The subdirectories of this directory
	Directories []*Directory
}

// ListAll queries for all files and directories recursively, starting at the
// given directory, for the given user. A `depth` can be provided to limit
// recursion, but only a value of `0` is supported (to disable recursion
// altogether).
func ListAll(ctx *data.Context, user *auth.User, depth int, dir *Directory) (*CompositeListing, error) {
	composite := &CompositeListing{Parent: dir}
	if depth == 0 {
		// If depth=0 then we only care about the directory, which at this
		// point should be the only field in the composite, so we can just
		// return early.
		return composite, nil
	}
	composite.Directories = ListDirectories(ctx, user, dir)
	composite.Files = ListFiles(ctx, user, dir)
	return composite, nil
}

// ListFiles will list all Files under the given Directory for the given user.
func ListFiles(ctx *data.Context, user *auth.User, dir *Directory) []*File {
	var entries []*File
	if dir == nil {
		ctx.DB.Where("user_id = ? and parent_id is null", user.ID).Find(&entries)
	} else {
		ctx.DB.Where("user_id = ? and parent_id = ?", user.ID, dir.ID).Preload("Parent").Find(&entries)
	}
	return entries
}

// ListDirectories will list all Directories under the given Directory for the
// given user.
func ListDirectories(ctx *data.Context, user *auth.User, dir *Directory) []*Directory {
	var entries []*Directory
	if dir == nil {
		ctx.DB.Where("user_id = ? and parent_id is null", user.ID).Find(&entries)
	} else {
		ctx.DB.Where("user_id = ? and parent_id = ?", user.ID, dir.ID).Preload("Parent").Find(&entries)
	}
	return entries
}

// FindFileByPath will find a single file for the given user at the given path.
func FindFileByPath(ctx *data.Context, user *auth.User, path string) (*File, error) {
	var file *File
	err := ctx.DB.Where("user_id = ? and full_name = ?", user.ID, path).First(&file).Error
	return file, err
}

// DeleteFile will remove the file from the index.
func DeleteFile(ctx *data.Context, user *auth.User, file *File) error {
	return ctx.DB.Where("user_id = ? and id = ?", user.ID, file.ID).Delete(file).Error
}

// DeleteDirectory will remove the given Directory from the index.
func DeleteDirectory(ctx *data.Context, user *auth.User, dir *Directory) error {
	return ctx.DB.Where("user_id = ? and id = ?", user.ID, dir.ID).Delete(dir).Error
}

// DeletePath will remove all files and directories recursively from the index,
// starting at the given path.
func DeletePath(ctx *data.Context, user *auth.User, path string) error {
	var entity interface{}
	var isDir bool = false
	entity, err := FindFileByPath(ctx, user, path)
	if err != nil {
		entity, err = FindDirByPath(ctx, user, path)
		if err != nil {
			return err
		}
		isDir = true
	}
	if isDir {
		return DeleteDirectory(ctx, user, entity.(*Directory))
	}
	return DeleteFile(ctx, user, entity.(*File))
}

// RenameFile will change the name of the file in the index to the given name.
func RenameFile(ctx *data.Context, file *File, name string) error {
	file.SetNames(name)
	return ctx.DB.Save(file).Error
}

// RenameDirectory will change the name of the directory and update all of its
// descendants in the index to keep their `FullName` values accurate.
func RenameDirectory(ctx *data.Context, user *auth.User, dir *Directory, name string) error {
	dir.SetNames(name)
	err := ctx.DB.Save(dir).Error
	if err != nil {
		return err
	}
	composite, err := ListAll(ctx, user, 1, dir)
	if err != nil {
		return err
	}
	for _, file := range composite.Files {
		// Files don't get "renamed", but this will trigger the reprojection of
		// the FullName, which will store the new parent name in it.
		err = RenameFile(ctx, file, file.Name)
		if err != nil {
			return err
		}
	}
	for _, subdir := range composite.Directories {
		// Subdirectories don't get "renamed", but this will trigger the
		// reprojection of the FullName, which will store the new parent name.
		err = RenameDirectory(ctx, user, subdir, subdir.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

// RenamePath will rename either a file or directory at the given path. See
// `RenameFile` and `RenameDirectory` for the individual implementations.
func RenamePath(ctx *data.Context, user *auth.User, src string, dest string) error {
	var entity interface{}
	var isDir bool = false
	entity, err := FindFileByPath(ctx, user, src)
	if err != nil {
		entity, err = FindDirByPath(ctx, user, src)
		if err != nil {
			return err
		}
		isDir = true
	}
	name := filepath.Base(dest)
	if isDir {
		return RenameDirectory(ctx, user, entity.(*Directory), name)
	}
	return RenameFile(ctx, entity.(*File), name)
}

// FindDirByPath will find a directory whose `Fullname` matches the given path.
func FindDirByPath(ctx *data.Context, user *auth.User, path string) (*Directory, error) {
	var file *Directory
	err := ctx.DB.Where("user_id = ? and full_name = ?", user.ID, path).First(&file).Error
	if err != nil {
		return &Directory{}, err
	}
	return file, nil
}

// CreateFile will insert the given File into the index.
func CreateFile(ctx *data.Context, file *File) (*File, error) {
	err := ctx.DB.Create(&file).Error
	return file, err
}

// CreateDir will insert the given Directory into the index.
func CreateDir(ctx *data.Context, directory *Directory) (*Directory, error) {
	err := ctx.DB.Create(directory).Error
	return directory, err
}
