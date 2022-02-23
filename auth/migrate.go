package auth

import (
	"go.b8s.dev/nucleus/data"
)

func AutoMigrate(c *data.Context) {
	err := c.DB.AutoMigrate(&User{}, &Credential{})
	if err != nil {
		panic(err)
	}
}
