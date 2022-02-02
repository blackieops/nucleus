package files

import (
	"errors"

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
		ctx.DB.Where("user_id = ? and parent_id = ?", user.ID, dir.ID).Find(&entries)
	}
	return entries
}

func ListDirectories(ctx *data.Context, user *auth.User, dir *Directory) []*Directory {
	var entries []*Directory
	if dir == nil {
		ctx.DB.Where("user_id = ? and parent_id is null", user.ID).Find(&entries)
	} else {
		ctx.DB.Where("user_id = ? and parent_id = ?", user.ID, dir.ID).Find(&entries)
	}
	return entries
}

func FindFile(ctx *data.Context, id int) (*File, error) {
	var file *File
	ctx.DB.First(&file, id)
	if file == nil {
		return &File{}, errors.New("Could not find file.")
	}
	return file, nil
}

func FindFileByPath(ctx *data.Context, user *auth.User, path string) (*File, error) {
	var file *File
	ctx.DB.Where("user_id = ? and full_name = ?", user.ID, path).First(&file)
	if file == nil {
		return &File{}, errors.New("Could not find file.")
	}
	return file, nil
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
