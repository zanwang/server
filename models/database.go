package models

import (
	"database/sql"

	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/majimoe/server/config"
	_ "github.com/mxk/go-sqlite/sqlite3"
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
	switch conf.Database.Type {
	case "sqlite3":
		DB = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	case "mysql":
		DB = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{conf.Database.Engine, conf.Database.Encoding}}
	}

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
