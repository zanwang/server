package main

import (
  "net/http"
  "github.com/go-martini/martini"
  "log"
  "strconv"
)

func main() {
  m := martini.Classic()
  config := Config.Get("server").(map[interface{}]interface{})
  host := config["host"].(string)
  port := config["port"].(int)

  m.Use(martini.Static("public"))

  log.Printf("Listening at port %s:%d", host, port)
  log.Fatal(http.ListenAndServe(host + ":" + strconv.Itoa(port), m))
}