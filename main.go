package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "log"
  "strconv"
  "github.com/tommy351/maji.moe/controllers"
)

func main() {
  r := mux.NewRouter()
  config := Config.Get("server").(map[interface{}]interface{})
  host := config["host"].(string)
  port := config["port"].(int)

  // Routes
  addRoute(r, "/", new(controllers.Home))
  addRoute(r, "/login", new(controllers.Login))
  addRoute(r, "/signup", new(controllers.Signup))

  // Serve static files
  r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

  // Start server
  log.Printf("Listening at http://%s:%d", host, port)
  panic(http.ListenAndServe(host + ":" + strconv.Itoa(port), r))
}

func addRoute(r *mux.Router, path string, app controllers.ControllerInterface) {
  // Initialize the controller
  app.Prepare()
  app.Init()
  app.AddMethod("GET", app.Get)
  app.AddMethod("POST", app.Post)
  app.AddMethod("PUT", app.Put)
  app.AddMethod("DELETE", app.Delete)
  app.AddMethod("PATCH", app.Patch)
  app.AddMethod("HEAD", app.Head)

  // Register the route
  r.HandleFunc(path, app.Handle)
}