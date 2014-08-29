package models

import (
	"database/sql"

	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3" // import sqlite3
	"github.com/tommy351/maji.moe/config"
)

var DB *gorp.DbMap

func init() {
	conf := config.Config

	// Connect to the database
	db, err := sql.Open(conf.Database.Type, conf.Database.Path)

	if err != nil {
		panic(err)
	}

	// Build a gorp map
	DB = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// Add tables
	DB.AddTableWithName(User{}, "users").SetKeys(true, "id")
	DB.AddTableWithName(Domain{}, "domains").SetKeys(true, "id")
	DB.AddTableWithName(Record{}, "records").SetKeys(true, "id")
	DB.AddTableWithName(Token{}, "tokens").SetKeys(true, "id")

	// Create tables
	if err = DB.CreateTablesIfNotExists(); err != nil {
		panic(err)
	}
}
