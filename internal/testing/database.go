package testing

import (
	"errors"

	"go.b8s.dev/nucleus/config"
	"go.b8s.dev/nucleus/data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// WithData opens a connection to a test database, starts a transaction and runs
// the given function with the test transaction in the `data.Context`. Always
// rolls back the transaction after the function executes.
func WithData(block func(*data.Context)) {
	// TODO: fix brittle hardcoded path
	conf, err := config.LoadConfig("../config.test.yaml")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(conf.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Transaction(func(tx *gorm.DB) error {
		block(&data.Context{DB: tx})
		return errors.New("End of test!")
	})
}
