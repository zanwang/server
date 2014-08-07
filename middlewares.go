package main

import (
  "github.com/go-martini/martini"
  "github.com/martini-contrib/sessions"
)

type CurrentUser struct {
  User
  LoggedIn bool
}

func GetCurrentUser(c martini.Context, s sessions.Session) {
  var user CurrentUser

  err := DbMap.SelectOne(&user, "SELECT * FROM users WHERE Id=?", s.Get("UserId"))

  if err != nil {
    user.LoggedIn = false
    s.Delete("UserId")
  } else {
    user.LoggedIn = true
  }

  c.Map(user)
}