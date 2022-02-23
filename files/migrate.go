package files

import (
	"go.b8s.dev/nucleus/data"
)

func AutoMigrate(c *data.Context) {
	err := c.DB.AutoMigrate(&File{}, &Directory{})
	if err != nil {
		panic(err)
	}
}
