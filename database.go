package main

import (
  "database/sql"
  "github.com/coopernurse/gorp"
  _ "github.com/mattn/go-sqlite3"
)

type (
  User struct {
    Id int64
    Username string
    Password string
    Email string
    CreatedAt int64
    UpdatedAt int64
    Activated bool
    DisplayName string
    ActivatedToken string
    LoggedIn bool `db:"-"`
  }

  Domain struct {
    Id int64
    Name string
    CreatedAt int64
    UpdatedAt int64
    UserId int64
    Public bool
  }

  Record struct {
    Id int64
    Type string
    Subdomain string
    Destination string
    CreatedAt int64
    UpdatedAt int64
    DomainId int64
  }
)

func initDb() *gorp.DbMap {
  // Read configuration
  config := Config.Get("database").(map[interface{}]interface{})
  dbType := config["type"].(string)
  dbPath := config["path"].(string)

  // Connect to the database
  db, err := sql.Open(dbType, dbPath)
  checkErr(err, "Failed to connect to the database")

  // Construct a gorp map
  dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

  // Add tables
  dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
  dbmap.AddTableWithName(Domain{}, "domains").SetKeys(true, "Id")
  dbmap.AddTableWithName(Record{}, "records").SetKeys(true, "Id")

  // Create tables
  err = dbmap.CreateTablesIfNotExists()
  checkErr(err, "Failed to create tables")

  return dbmap
}