package data

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Context wraps the database backend.
type Context struct {
	DB *gorm.DB
}

// Connect establishes a database connection pool, connecting to the given URL.
func Connect(databaseUrl string) *Context {
	gormDB, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &Context{DB: gormDB}
}
