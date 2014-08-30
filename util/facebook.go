package util

import (
	"github.com/huandu/facebook"
	"github.com/tommy351/maji.moe/config"
)

var Facebook *facebook.App

func init() {
	conf := config.Config
	Facebook = facebook.New(conf.Facebook.AppID, conf.Facebook.AppSecret)
}
