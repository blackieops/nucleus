package auth

import (
	"com.blackieops.nucleus/data"
)

func AutoMigrate(c *data.Context) {
	c.DB.AutoMigrate(&User{})
	c.DB.AutoMigrate(&Credential{})
}
