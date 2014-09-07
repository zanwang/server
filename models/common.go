package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/majimoe/server/config"
)

const (
	ISOTimeFormat = "2006-01-02T15:04:05Z"
)

var DB gorm.DB

func init() {
	conf := config.Config
	db, err := gorm.Open(conf.Database.Type, conf.Database.Path)

	if err != nil {
		panic(err)
	}

	DB = db
}

func ISOTime(t time.Time) string {
	return t.UTC().Format(ISOTimeFormat)
}
