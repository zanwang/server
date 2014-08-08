package main

import (
  "net/http"
  "crypto/sha1"
  "io"
  "encoding/base64"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/sessions"
  "github.com/martini-contrib/render"
  "github.com/dchest/uniuri"
)

func GetCurrentUser(c martini.Context, s sessions.Session) {
  var user User

  err := DbMap.SelectOne(&user, "SELECT * FROM users WHERE Id=?", s.Get("UserId"))

  if err != nil {
    user.LoggedIn = false
    s.Delete("UserId")
  } else {
    user.LoggedIn = true
  }

  c.Map(user)
}

func NeedLogin(user User, r render.Render) {
  if !user.LoggedIn {
    r.Redirect("/login")
    return
  }
}

type CsrfToken struct {
  Hash string
  Secret string
  Token string
  session sessions.Session
}

func (c *CsrfToken) GetToken() string {
  if c.Token != "" {
    return c.Token
  }

  c.Token = csrfTokenize(c.Hash, c.Secret)
  c.session.Set("csrfSecret", c.Hash)

  return c.Token
}

func csrfTokenize(hash string, secret string) string {
  h := sha1.New()

  io.WriteString(h, hash)
  io.WriteString(h, secret)

  return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func CsrfMiddleware(c martini.Context, s sessions.Session, req *http.Request, config ServerConfig, res http.ResponseWriter) {
  token := CsrfToken{
    Hash: uniuri.New(),
    Secret: config.Secret,
    session: s,
  }

  c.Map(&token)

  switch req.Method {
  case "GET", "HEAD", "OPTIONS":
    return
  }

  hash := s.Get("csrfSecret").(string)

  if csrfTokenize(hash, config.Secret) != req.FormValue("_csrf") {
    http.Error(res, "Invalid csrf token.", http.StatusBadRequest)
  }
}