package auth

import (
	"testing"

	"com.blackieops.nucleus/data"
	testUtils "com.blackieops.nucleus/internal/testing"
)

func TestFindUserCredentials(t *testing.T) {
	testUtils.WithData(func(c *data.Context) {
		user, err := CreateUser(c, &User{
			Name: "Test",
			Username: "test",
			EmailAddress: "test@example.com",
		})
		if err != nil {
			t.Errorf("Failed to create user: %v", err)
		}
		wrongUser, err := CreateUser(c, &User{
			Name: "Wrong",
			Username: "wrong",
			EmailAddress: "wrong@example.com",
		})
		if err != nil {
			t.Errorf("Failed to create user: %v", err)
		}
		cred := &Credential{Data: "abc123"}
		_, err = CreateCredential(c, user, cred)
		if err != nil {
			t.Errorf("Failed to create credential: %v", err)
		}
		wrongCred := &Credential{Data: "wrong123"}
		_, err = CreateCredential(c, wrongUser, wrongCred)
		if err != nil {
			t.Errorf("Failed to create credential: %v", err)
		}

		found, err := FindUserCredentials(c, user)
		if err != nil {
			t.Errorf("Unexpected error in FindUserCredentials: %v", err)
		}
		if len(found) != 1 {
			t.Errorf("Unexpected result count for FindUserCredentials: %d", len(found))
		}
		if found[0].ID != cred.ID {
			t.Errorf("Found incorrect credential! %v", found[0].ID)
		}
	})
}

func TestCreateCredential(t *testing.T) {
	testUtils.WithData(func(c *data.Context) {
		user, err := CreateUser(c, &User{
			Name: "Test",
			Username: "test",
			EmailAddress: "test@example.com",
		})
		if err != nil {
			t.Errorf("Failed to create user: %v", err)
		}

		// Test basic creation
		valid := &Credential{Data: "abc123"}
		created, err := CreateCredential(c, user, valid)
		if err != nil {
			t.Errorf("Failed to create credential: %v", err)
		}

		// Test that credential's Data field is not persisted
		var found *Credential
		c.DB.Where("id = ?", created.ID).First(&found)
		if found.Data != "" {
			t.Errorf("Temporary raw data field was persisted!! Data: %v", created.Data)
		}

		// Test that invalid credentials return error
		invalid := &Credential{}
		_, err = CreateCredential(c, user, invalid)
		if err == nil {
			t.Errorf("Invalid credential was created!")
		}
	})
}

func TestValidateCredential(t *testing.T) {
	testUtils.WithData(func(c *data.Context) {
		user, err := CreateUser(c, &User{
			Name: "Test",
			Username: "test",
			EmailAddress: "test@example.com",
		})
		if err != nil {
			t.Errorf("Failed to create user: %v", err)
		}
		fixture := &Credential{Data: "password123", Type: CredentialTypePassword}
		cred, err := CreateCredential(c, user, fixture)
		if err != nil {
			t.Errorf("Failed to create credential: %v", err)
		}
		err = ValidateCredential(cred, "password123")
		if err != nil {
			t.Errorf("Failed to validate credential: %v", err)
		}
	})
}

func TestDeleteCredential(t *testing.T) {
	testUtils.WithData(func(c *data.Context) {
		user, err := CreateUser(c, &User{
			Name: "Test",
			Username: "test",
			EmailAddress: "test@example.com",
		})
		if err != nil {
			t.Errorf("Failed to create user: %v", err)
		}
		fixture := &Credential{Data: "password123", Type: CredentialTypePassword}
		cred, err := CreateCredential(c, user, fixture)
		if err != nil {
			t.Errorf("Failed to create credential: %v", err)
		}

		err = DeleteCredential(c, cred)
		if err != nil {
			t.Errorf("Failed to delete credential: %v", err)
		}

		found, err := FindUserCredentials(c, user)
		if err != nil {
			t.Errorf("Failed to list user credentials: %v", err)
		}
		if len(found) != 0 {
			t.Errorf("Found unexpected credentials after deletion!")
		}
	})
}
