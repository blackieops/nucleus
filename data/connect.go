package data

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Context struct {
	DB *gorm.DB
}

func Connect(databaseUrl string) *Context {
	gormDB, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &Context{DB: gormDB}
}
