package auth

import (
	"go.b8s.dev/nucleus/data"
)

// AutoMigrate runs database migrations for this package's models.
func AutoMigrate(c *data.Context) {
	err := c.DB.AutoMigrate(&User{}, &Credential{})
	if err != nil {
		panic(err)
	}
}
