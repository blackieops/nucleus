package nxc

import (
	"go.b8s.dev/nucleus/data"
)

func AutoMigrate(c *data.Context) {
	err := c.DB.AutoMigrate(&NextcloudAppPassword{}, &NextcloudAuthSession{})
	if err != nil {
		panic(err)
	}
}
