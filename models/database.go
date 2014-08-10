package models

import (
	"database/sql"

	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3" // import sqlite3
	"github.com/tommy351/maji.moe/config"
)

// Load initializes the database
func Load(conf *config.Config) *gorp.DbMap {
	dbType := conf.Database.Type
	dbPath := conf.Database.Path

	// Connect to the database
	db, err := sql.Open(dbType, dbPath)

	if err != nil {
		panic(err)
	}

	// Build a gorp map
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// Add tables
	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "id")
	dbMap.AddTableWithName(Domain{}, "domains").SetKeys(true, "id")
	dbMap.AddTableWithName(Record{}, "records").SetKeys(true, "id")
	dbMap.AddTableWithName(Token{}, "tokens").SetKeys(true, "id")

	// Create tables
	if err = dbMap.CreateTablesIfNotExists(); err != nil {
		panic(err)
	}

	return dbMap
}
