package controllers

type Home struct {
  Controller
}

func (c *Home) Get() {
  c.Render("index", nil)
}