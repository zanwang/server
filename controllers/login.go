package controllers

type Login struct {
  Controller
}

func (c *Login) Get() {
  c.Render("login", nil)
}

func (c *Login) Post() {
  //
}