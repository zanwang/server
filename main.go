package main

import (
	"strconv"

	"github.com/majimoe/server/config"
	"github.com/majimoe/server/models"
	"github.com/majimoe/server/server"
)

func main() {
	conf := config.Config
	host := conf.Server.Host
	port := conf.Server.Port

	server.Server().Run(host + ":" + strconv.Itoa(port))
	defer models.DB.Close()
}
