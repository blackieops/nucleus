package files

import (
	"testing"

	"go.b8s.dev/nucleus/auth"
	"go.b8s.dev/nucleus/data"
	testUtils "go.b8s.dev/nucleus/internal/testing"
)

func TestDeleteFile(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, err := auth.CreateUser(ctx, user)
		if err != nil {
			t.Errorf("Failed to setup test user: %v", err)
		}
		file, err := CreateFile(ctx, &File{Name: "butt.txt", FullName: "butt.txt", User: *user})
		if err != nil {
			t.Errorf("Failed to setup test file: %v", err)
		}
		err = DeleteFile(ctx, testUser, file)
		if err != nil {
			t.Errorf("Failed to delete file: %v", err)
		}
	})
}

func TestDeleteDirectory(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, err := auth.CreateUser(ctx, user)
		if err != nil {
			t.Errorf("Failed to setup test user: %v", err)
		}
		dir, err := CreateDir(ctx, &Directory{Name: "things", FullName: "things", User: *user})
		if err != nil {
			t.Errorf("Failed to setup test directory: %v", err)
		}
		subdir, err := CreateDir(ctx, &Directory{Name: "more things", FullName: "things/more things", User: *user, Parent: dir})
		if err != nil {
			t.Errorf("Failed to setup test sub-directory: %v", err)
		}
		file, err := CreateFile(ctx, &File{Name: "butt.txt", FullName: "things/butt.txt", User: *user, Parent: dir})
		if err != nil {
			t.Errorf("Failed to setup test file: %v", err)
		}
		err = DeleteDirectory(ctx, testUser, dir)
		if err != nil {
			t.Errorf("Failed to delete directory: %v", err)
		}
		err = ctx.DB.Where("id = ?", file.ID).First(&file).Error
		if err == nil {
			t.Errorf("Directory did not delete its files!")
		}
		err = ctx.DB.Where("id = ?", subdir.ID).First(&subdir).Error
		if err == nil {
			t.Errorf("Directory did not delete its subdirectories!")
		}
	})
}

func TestDeletePathWhenFile(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, err := auth.CreateUser(ctx, user)
		if err != nil {
			t.Errorf("Failed to setup test user: %v", err)
		}
		file, err := CreateFile(ctx, &File{Name: "butt.txt", FullName: "butt.txt", User: *user})
		if err != nil {
			t.Errorf("Failed to setup test file: %v", err)
		}

		err = DeletePath(ctx, user, "butt.txt")
		if err != nil {
			t.Errorf("Failed to delete file by path: %v", err)
		}

		err = ctx.DB.Where("id = ?", file.ID).First(&file).Error
		if err == nil {
			t.Errorf("File did not actually get deleted!")
		}
	})
}

func TestDeletePathWhenDirectory(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, err := auth.CreateUser(ctx, user)
		if err != nil {
			t.Errorf("Failed to setup test user: %v", err)
		}
		dir, err := CreateDir(ctx, &Directory{Name: "things", FullName: "things", User: *user})
		if err != nil {
			t.Errorf("Failed to setup test directory: %v", err)
		}
		file, err := CreateFile(ctx, &File{Name: "butt.txt", FullName: "things/butt.txt", User: *user, Parent: dir})
		if err != nil {
			t.Errorf("Failed to setup test file: %v", err)
		}

		err = DeletePath(ctx, user, "things")
		if err != nil {
			t.Errorf("Failed to delete directory by path: %v", err)
		}

		err = ctx.DB.Where("id = ?", file.ID).First(&file).Error
		if err == nil {
			t.Errorf("Directory files did not get deleted!")
		}
	})
}

func TestRenamePathWhenDirectory(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, err := auth.CreateUser(ctx, user)
		if err != nil {
			t.Errorf("Failed to setup test user: %v", err)
		}
		dir, err := CreateDir(ctx, &Directory{Name: "things", FullName: "things", User: *user})
		if err != nil {
			t.Errorf("Failed to setup test directory: %v", err)
		}
		_, err = CreateFile(ctx, &File{Name: "butt.txt", FullName: "things/butt.txt", User: *user, Parent: dir})
		if err != nil {
			t.Errorf("Failed to setup test file: %v", err)
		}
		err = RenamePath(ctx, user, "things", "stuff")
		if err != nil {
			t.Errorf("Failed to rename directory by path: %v", err)
		}
		err = ctx.DB.Where("full_name = ?", "things/butt.txt").First(&File{}).Error
		if err == nil {
			t.Errorf("Directory's File's FullName was not updated!")
		}
	})
}

func TestRenamePathWhenFile(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, err := auth.CreateUser(ctx, user)
		if err != nil {
			t.Errorf("Failed to setup test user: %v", err)
		}

		// When the file is in a parent directory
		dir, err := CreateDir(ctx, &Directory{Name: "things", FullName: "things", User: *user})
		if err != nil {
			t.Errorf("Failed to setup test directory: %v", err)
		}
		file, err := CreateFile(ctx, &File{Name: "butt.txt", FullName: "things/butt.txt", User: *user, Parent: dir})
		if err != nil {
			t.Errorf("Failed to setup test file: %v", err)
		}

		err = RenamePath(ctx, user, "things/butt.txt", "things/hello.txt")
		if err != nil {
			t.Errorf("Failed to rename file by path: %v", err)
		}
		err = ctx.DB.Where("full_name = ?", "things/butt.txt").First(&File{}).Error
		if err == nil {
			t.Errorf("File's FullName was not updated!")
		}

		// When the file has no parent directory
		file, err = CreateFile(ctx, &File{Name: "alone.txt", FullName: "alone.txt", User: *user, Parent: nil})
		if err != nil {
			t.Errorf("Failed to setup test file: %v", err)
		}
		err = RenamePath(ctx, user, "alone.txt", "sup.txt")
		if err != nil {
			t.Errorf("Failed to rename file by path: %v", err)
		}
		err = ctx.DB.Where("full_name = ?", "alone.txt").First(&file).Error
		if err == nil {
			t.Errorf("File's FullName was not updated!")
		}
	})
}

func TestRenameDirectory(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		// Set up user
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, err := auth.CreateUser(ctx, user)
		if err != nil {
			t.Errorf("Failed to setup test user: %v", err)
			return
		}

		// Set up parent directory
		dir, err := CreateDir(ctx, &Directory{Name: "things", FullName: "things", User: *user})
		if err != nil {
			t.Errorf("Failed to setup test directory: %v", err)
			return
		}
		subdir, err := CreateDir(ctx, &Directory{Name: "important", FullName: "things/important", Parent: dir, User: *user})
		if err != nil {
			t.Errorf("Failed to setup test directory: %v", err)
			return
		}
		_, err = CreateFile(ctx, &File{Name: "child.jpg", FullName: "things/important/child.jpg", User: *user, Parent: subdir})
		if err != nil {
			t.Errorf("Failed to setup test file: %v", err)
			return
		}

		// Try renaming it
		err = RenameDirectory(ctx, user, dir, "butt")
		if err != nil {
			t.Errorf("Failed to rename an empty directory: %v", err)
		}
		_, err = FindDirByPath(ctx, user, "butt")
		if err != nil {
			t.Errorf("Failed to find directory after rename: %v", err)
		}
		_, err = FindDirByPath(ctx, user, "butt/important")
		if err != nil {
			t.Errorf("Failed to find sub-directory after rename: %v", err)
		}
		_, err = FindFileByPath(ctx, user, "butt/important/child.jpg")
		if err != nil {
			t.Errorf("Failed to find sub-directory's file after rename: %v", err)
		}
	})
}

func TestRenameFile(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		// Set up user
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, err := auth.CreateUser(ctx, user)
		if err != nil {
			t.Errorf("Failed to setup test user: %v", err)
			return
		}

		// Set up file
		dir, err := CreateDir(ctx, &Directory{Name: "Desktop", FullName: "Desktop", User: *user})
		if err != nil {
			t.Errorf("Failed to setup test directory: %v", err)
		}
		file, err := CreateFile(ctx, &File{Name: "test.docx", FullName: "Desktop/test.docx", User: *user, Parent: dir})
		if err != nil {
			t.Errorf("Failed to setup test file: %v", err)
		}

		// Test that it renames
		err = RenameFile(ctx, file, "Proposal.docx")
		if err != nil {
			t.Errorf("Failed to rename file: %v", err)
		}

		_, err = FindFileByPath(ctx, user, "Desktop/Proposal.docx")
		if err != nil {
			t.Errorf("Failed to find file after rename: %v", err)
		}
	})
}
