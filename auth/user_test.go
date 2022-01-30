package auth

import (
	"testing"

	"com.blackieops.nucleus/data"
	testUtils "com.blackieops.nucleus/internal/testing"
)

func TestFindUserByUsername(t *testing.T) {
	testUtils.WithData(func(ctx *data.Context) {
		AutoMigrate(ctx)
		ctx.DB.Create(&User{Username: "admin", Name: "Admin", EmailAddress: "A@example.com"})

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
