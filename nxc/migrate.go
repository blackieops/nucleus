package nxc

import (
	"com.blackieops.nucleus/data"
)

func AutoMigrate(c *data.Context) {
	err := c.DB.AutoMigrate(&NextcloudAppPassword{}, &NextcloudAuthSession{})
	if err != nil {
		panic(err)
	}
}
