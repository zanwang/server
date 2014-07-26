package main

import (
  "net/http"
  "github.com/go-martini/martini"
  "log"
  "strconv"
)

func main() {
  server := martini.Classic()
  port := Config.Get("server").(map[interface{}]interface{})["port"].(int)

  log.Printf("Listening at port %d", port)
  log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), server))
}