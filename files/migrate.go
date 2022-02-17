package files

import (
	"com.blackieops.nucleus/data"
)

func AutoMigrate(c *data.Context) {
	err := c.DB.AutoMigrate(&File{}, &Directory{})
	if err != nil {
		panic(err)
	}
}
