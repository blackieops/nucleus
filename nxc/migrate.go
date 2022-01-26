package nxc

import (
	"com.blackieops.nucleus/data"
)

func AutoMigrate(c *data.Context) {
	c.DB.AutoMigrate(&NextcloudAppPassword{}, &NextcloudAuthSession{})
}
