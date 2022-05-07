package nxc

import (
	"testing"

	"go.b8s.dev/nucleus/auth"
	"go.b8s.dev/nucleus/data"
	testUtils "go.b8s.dev/nucleus/internal/testing"
)

func TestListNextcloudAppPasswordsForUser(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, _ = auth.CreateUser(ctx, user)
		wrongUser := &auth.User{Name: "Wrong", Username: "wrong", EmailAddress: "wrong@example.com"}
		wrongUser, _ = auth.CreateUser(ctx, wrongUser)
		session, _ := CreateNextcloudAuthSession(ctx)
		password, _ := CreateNextcloudAppPassword(ctx, session, user)
		CreateNextcloudAppPassword(ctx, session, wrongUser)

		result, err := ListNextcloudAppPasswordsForUser(ctx, user)

		if err != nil {
			t.Errorf("ListNextcloudAppPasswordsForUser had unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("ListNextcloudAppPasswordsForUser had unexpected result length: %v", len(result))
		}
		if result[0].ID != password.ID {
			t.Errorf("ListNextcloudAppPasswordsForUser returned unexpected password: %v", result[0].ID)
		}
	})
}

func TestFindNexcloudAppPassword(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, _ = auth.CreateUser(ctx, user)
		wrongUser := &auth.User{Name: "Wrong", Username: "wrong", EmailAddress: "wrong@example.com"}
		wrongUser, _ = auth.CreateUser(ctx, wrongUser)
		session, _ := CreateNextcloudAuthSession(ctx)
		password, _ := CreateNextcloudAppPassword(ctx, session, user)
		wrongPassword, _ := CreateNextcloudAppPassword(ctx, session, wrongUser)

		result, err := FindNextcloudAppPassword(ctx, user, password.ID)
		if err != nil {
			t.Errorf("FindNextcloudAppPassword had unexpected error: %v", err)
		}
		if result.ID != password.ID {
			t.Errorf("FindNextcloudAppPassword found wrong password: %v", result.ID)
		}

		_, err = FindNextcloudAppPassword(ctx, user, wrongPassword.ID)
		if err == nil {
			t.Errorf("FindNextcloudAppPassword found a password that it shouldn't've!")
		}
	})
}

func TestCreateNextcloudAppPassword(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, _ = auth.CreateUser(ctx, user)
		session, _ := CreateNextcloudAuthSession(ctx)

		result, err := CreateNextcloudAppPassword(ctx, session, user)

		if err != nil {
			t.Errorf("CreateNextcloudAppPassword encountered unexpected error: %v", err)
		}
		if result.PasswordDigest == "" {
			t.Errorf("CreateNextcloudAppPassword did not generate a password!")
		}
		if uint(result.UserID) != user.ID {
			t.Errorf("CreateNextcloudAppPassword did not set the right user ID: %v", result.UserID)
		}
	})
}

func TestFindNextcloudAppPasswordByPassword(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, _ = auth.CreateUser(ctx, user)
		session, _ := CreateNextcloudAuthSession(ctx)
		password, _ := CreateNextcloudAppPassword(ctx, session, user)

		result, err := FindNextcloudAppPasswordByPassword(ctx, session.RawAppPassword)

		if err != nil {
			t.Errorf("FindNextcloudAppPasswordByPassword encountered unexpected error: %v", err)
		}
		if result.ID != password.ID {
			t.Errorf("FindNextcloudAppPasswordByPassword did not find the right password: %v", result.ID)
		}
	})
}

func TestDeleteNextcloudAppPassword(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, _ = auth.CreateUser(ctx, user)
		session, _ := CreateNextcloudAuthSession(ctx)
		password, _ := CreateNextcloudAppPassword(ctx, session, user)

		err := DeleteNextcloudAppPassword(ctx, password)

		if err != nil {
			t.Errorf("NextcloudAppPassword was not deleted! %v", err)
		}

		_, err = FindNextcloudAppPassword(ctx, user, password.ID)
		if err == nil {
			t.Errorf("Should not have found password after deletion!")
		}
	})
}
