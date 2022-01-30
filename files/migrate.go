package files

import (
	"com.blackieops.nucleus/data"
)

func AutoMigrate(c *data.Context) {
	c.DB.AutoMigrate(&File{}, &Directory{})
}
