package files

import (
	"testing"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/data"
	testUtils "com.blackieops.nucleus/internal/testing"
)

func TestDeleteFile(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		AutoMigrate(ctx)
		auth.AutoMigrate(ctx)
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
		AutoMigrate(ctx)
		auth.AutoMigrate(ctx)
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
		err = DeleteDirectory(ctx, testUser, dir)
		if err != nil {
			t.Errorf("Failed to delete directory: %v", err)
		}
		err = ctx.DB.Where("id = ?", file.ID).First(&file).Error
		if err == nil {
			t.Errorf("Directory did not delete its files!")
		}
	})
}
