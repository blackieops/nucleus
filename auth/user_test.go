package auth

import (
	"testing"

	"com.blackieops.nucleus/data"
	testUtils "com.blackieops.nucleus/internal/testing"
)

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
