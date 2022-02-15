package files

import (
	"path/filepath"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/data"
)

type CompositeListing struct {
	// Parent is the directory for under which these files and directories sit
	Parent *Directory

	// The files in this directory
	Files []*File

	// The subdirectories of this directory
	Directories []*Directory
}

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

func ListFiles(ctx *data.Context, user *auth.User, dir *Directory) []*File {
	var entries []*File
	if dir == nil {
		ctx.DB.Where("user_id = ? and parent_id is null", user.ID).Find(&entries)
	} else {
		ctx.DB.Where("user_id = ? and parent_id = ?", user.ID, dir.ID).Preload("Parent").Find(&entries)
	}
	return entries
}

func ListDirectories(ctx *data.Context, user *auth.User, dir *Directory) []*Directory {
	var entries []*Directory
	if dir == nil {
		ctx.DB.Where("user_id = ? and parent_id is null", user.ID).Find(&entries)
	} else {
		ctx.DB.Where("user_id = ? and parent_id = ?", user.ID, dir.ID).Preload("Parent").Find(&entries)
	}
	return entries
}

func FindFile(ctx *data.Context, id int) (*File, error) {
	var file *File
	err := ctx.DB.First(&file, id).Error
	return file, err
}

func FindFileByPath(ctx *data.Context, user *auth.User, path string) (*File, error) {
	var file *File
	err := ctx.DB.Where("user_id = ? and full_name = ?", user.ID, path).First(&file).Error
	return file, err
}

func DeleteFile(ctx *data.Context, user *auth.User, file *File) error {
	return ctx.DB.Where("id = ?", file.ID).Delete(file).Error
}

func DeleteDirectory(ctx *data.Context, user *auth.User, dir *Directory) error {
	return ctx.DB.Where("id = ?", dir.ID).Delete(dir).Error
}

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

func RenameFile(ctx *data.Context, user *auth.User, file *File, name string) error {
	file.SetNames(name)
	return ctx.DB.Save(file).Error
}

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
		err = RenameFile(ctx, user, file, file.Name)
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
	return RenameFile(ctx, user, entity.(*File), name)
}

func FindDirByPath(ctx *data.Context, user *auth.User, path string) (*Directory, error) {
	var file *Directory
	err := ctx.DB.Where("user_id = ? and full_name = ?", user.ID, path).First(&file).Error
	if err != nil {
		return &Directory{}, err
	}
	return file, nil
}

func CreateFile(ctx *data.Context, file *File) (*File, error) {
	err := ctx.DB.Create(&file).Error
	return file, err
}

func CreateDir(ctx *data.Context, directory *Directory) (*Directory, error) {
	err := ctx.DB.Create(directory).Error
	return directory, err
}
