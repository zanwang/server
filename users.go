package main

import (
  "net/http"
  "time"
  "regexp"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/binding"
  "code.google.com/p/go.crypto/bcrypt"
  "github.com/martini-contrib/sessions"
  "github.com/dchest/uniuri"
)

var (
  rEmail = regexp.MustCompile(".+@.+\\..+")
)

type UserCreateForm struct {
  Username string `form:"username" binding:"required"`
  Password string `form:"password" binding:"required"`
  Confirm string `form:"confirm" binding:"required"`
  Email string `form:"email" binding:"required"`
}

func (form UserCreateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
  usernameLen := len(form.Username)
  if usernameLen > 20 || usernameLen < 3 {
    errors.Add([]string{"username"}, "FormatError", "The length of username must between 3 ~ 20.")
  }

  passwordLen := len(form.Password)
  if passwordLen > 50 || passwordLen < 6 {
    errors.Add([]string{"password"}, "FormatError", "The length of password must between 6 ~ 50.")
  }

  if form.Password != form.Confirm {
    errors.Add([]string{"confirm"}, "ContentError", "Password confirmation doesn't match.")
  }

  if !rEmail.MatchString(form.Email) {
    errors.Add([]string{"email"}, "ContentError", "Email is invalid.")
  }

  return errors
}

func UserNew(r render.Render, currentUser CurrentUser) {
  if currentUser.LoggedIn {
    r.Redirect("/users/" + currentUser.Username)
    return
  }

  r.HTML(http.StatusOK, "users/new", nil)
}

func UserCreate(form UserCreateForm, r render.Render, errors binding.Errors, s sessions.Session) {
  if errors != nil {
    r.HTML(http.StatusBadRequest, "users/new", map[string]interface{}{
      "Errors": formatErr(errors),
    })
    return
  }

  // Check whether user is registered
  var user User
  err := DbMap.SelectOne(&user, "SELECT * FROM users WHERE Username=? OR Email=?", form.Username, form.Email)

  if err == nil {
    errors = binding.Errors{}

    if form.Username == user.Username {
      errors.Add([]string{"username"}, "ContentError", "Username has been used.")
    }

    if form.Email == user.Email {
      errors.Add([]string{"email"}, "ContentError", "Email has been used.")
    }

    r.HTML(http.StatusBadRequest, "users/new", map[string]interface{}{
      "Errors": formatErr(errors),
    })
    return
  }

  password, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
  if err != nil {
    r.HTML(http.StatusInternalServerError, "errors/500", err)
    return
  }

  token := uniuri.New()
  now := time.Now().UnixNano()

  user = User{
    Username: form.Username,
    Password: string(password),
    Email: form.Email,
    CreatedAt: now,
    UpdatedAt: now,
    Activated: false,
    DisplayName: form.Username,
    ActivatedToken: token,
  }

  err = DbMap.Insert(&user)
  if err != nil {
    r.HTML(http.StatusInternalServerError, "errors/500", err)
    return
  }

  s.Set("UserId", user.Id)
  r.Redirect("/users/" + user.Username)
}

func UserShow(r render.Render, params martini.Params, currentUser CurrentUser) {
  var user User
  err := DbMap.SelectOne(&user, "SELECT * FROM users WHERE Username=?", params["user_id"])

  if err != nil {
    r.HTML(http.StatusNotFound, "errors/404", err)
    return
  }

  r.HTML(http.StatusOK, "users/show", map[string]interface{}{
    "User": user,
    "Editable": currentUser.Id == user.Id,
  })
}

func UserEdit(r render.Render, currentUser CurrentUser) {
  if !currentUser.LoggedIn {
    r.Redirect("/login")
    return
  }

  r.HTML(http.StatusOK, "users/edit", map[string]interface{}{
    "User": currentUser,
  })
}

func UserUpdate() {
  //
}

func UserDestroy() {
  //
}

func UserConfirm(r render.Render) {
  // r.Redirect("/users/")
}