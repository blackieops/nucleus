package auth

import (
	"testing"

	"com.blackieops.nucleus/data"
	testUtils "com.blackieops.nucleus/internal/testing"
)

func TestFindAllUsers(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		CreateUser(ctx, &User{Username: "admin", Name: "Admin", EmailAddress: "A@example.com"})
		CreateUser(ctx, &User{Username: "alice", Name: "Alice", EmailAddress: "B@example.com"})
		users := FindAllUsers(ctx)
		if len(users) != 2 {
			t.Errorf("Found incorrect number of users: %d", len(users))
		}
		for _, u := range users {
			if u.Username != "admin" && u.Username != "alice" {
				t.Errorf("FindAllUsers found unexpected user: %v", u.Username)
			}
		}
	})
}

func TestFindUser(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user, err := CreateUser(ctx, &User{Username: "admin", Name: "Admin", EmailAddress: "A@example.com"})
		if err != nil {
			t.Errorf("Failed to set up user: %v", err)
		}
		found := FindUser(ctx, int(user.ID))
		if found.ID != user.ID {
			t.Errorf("Found incorrect user: %d instead of %d", found.ID, user.ID)
		}
	})
}

func TestFindUserByUsername(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		CreateUser(ctx, &User{Username: "admin", Name: "Admin", EmailAddress: "A@example.com"})

		user, err := FindUserByUsername(ctx, "admin")
		if err != nil {
			t.Errorf("FindUserByUsername encountered a query error: %v", err)
		}
		if user.Username != "admin" {
			t.Errorf("FindUserByUsername found wrong user: %s", user.Username)
		}

		_, err = FindUserByUsername(ctx, "bobbytables")
		if err == nil {
			t.Errorf("FindUserByUsername found user when it shouldn't have.")
		}
	})
}

func TestCreateUser(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		user, err := CreateUser(ctx, &User{
			Name: "Tester",
			Username: "test",
			EmailAddress: "test@example.com",
		})
		if err != nil {
			t.Errorf("Failed to persist user: %v", err)
		}
		if user.ID == uint(0) {
			t.Errorf("User did not get created with a valid ID!")
		}
	})
}
