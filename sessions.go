package main

import (
  "net/http"
  "github.com/martini-contrib/render"
  "code.google.com/p/go.crypto/bcrypt"
  "github.com/martini-contrib/sessions"
  "github.com/martini-contrib/binding"
)

func SessionNew(r render.Render, user User, csrf *CsrfToken) {
  if user.LoggedIn {
    r.Redirect("/users/" + user.Username)
    return
  }

  r.HTML(http.StatusOK, "sessions/new", map[string]interface{}{
    "Token": csrf.GetToken(),
  })
}

type SessionCreateForm struct {
  Login string `form:"login" binding:"required"`
  Password string `form:"password" binding:"required"`
}

func SessionCreate(form SessionCreateForm, r render.Render, s sessions.Session, errors binding.Errors, csrf *CsrfToken) {
  if errors != nil {
    r.HTML(http.StatusBadRequest, "sessions/new", map[string]interface{}{
      "Errors": formatErr(errors),
      "Token": csrf.GetToken(),
    })
    return
  }

  var user User
  err := DbMap.SelectOne(&user, "SELECT * FROM users WHERE Username=? OR Email=?", form.Login, form.Login)

  if err != nil {
    errors = binding.Errors{}
    errors.Add([]string{"login"}, "ContentError", "User does not exist.")

    r.HTML(http.StatusUnauthorized, "sessions/new", map[string]interface{}{
      "Errors": formatErr(errors),
      "Token": csrf.GetToken(),
    })
    return
  }

  err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))

  if err == nil {
    s.Set("UserId", user.Id)
    r.Redirect("/users/" + user.Username)
  } else {
    errors = binding.Errors{}

    errors.Add([]string{"password"}, "ContentError", "Password does not match.")
    r.HTML(http.StatusUnauthorized, "sessions/new", map[string]interface{}{
      "Errors": formatErr(errors),
      "Token": csrf.GetToken(),
    })
  }
}

func SessionDestroy(s sessions.Session, r render.Render) {
  s.Delete("UserId")
  r.Redirect("/login")
}