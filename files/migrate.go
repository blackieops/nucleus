package files

import (
	"go.b8s.dev/nucleus/data"
)

// AutoMigrate runs database migrations for the models in this package.
func AutoMigrate(c *data.Context) {
	err := c.DB.AutoMigrate(&File{}, &Directory{})
	if err != nil {
		panic(err)
	}
}
