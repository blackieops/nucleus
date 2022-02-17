package auth

import (
	"com.blackieops.nucleus/data"
)

func AutoMigrate(c *data.Context) {
	err := c.DB.AutoMigrate(&User{}, &Credential{})
	if err != nil {
		panic(err)
	}
}
