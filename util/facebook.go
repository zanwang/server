package util

import (
	"github.com/huandu/facebook"
	"github.com/majimoe/server/config"
)

var Facebook *facebook.App

func init() {
	conf := config.Config
	Facebook = facebook.New(conf.Facebook.AppId, conf.Facebook.AppSecret)
}
