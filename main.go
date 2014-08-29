package main

import (
	"strconv"

	"github.com/tommy351/maji.moe/config"
	"github.com/tommy351/maji.moe/server"
)

func main() {
	conf := config.Config
	host := conf.Server.Host
	port := conf.Server.Port

	server.Server().Run(host + ":" + strconv.Itoa(port))
}
