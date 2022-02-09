package main

import (
	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/data"
)

func seedData(c *data.Context) {
	user, err := auth.CreateUser(c, &auth.User{Username: "admin", Name: "Admin", EmailAddress: "admin@example.com"})
	maybePanic(err)
	_, err = auth.CreateCredential(c, user, &auth.Credential{Type: auth.CredentialTypePassword, Data: "password123"})
	maybePanic(err)
}

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}
