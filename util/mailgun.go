package util

import (
	"github.com/mailgun/mailgun-go"
	"github.com/majimoe/server/config"
)

var Mailgun mailgun.Mailgun

func init() {
	conf := config.Config
	Mailgun = mailgun.NewMailgun(conf.Mailgun.Domain, conf.Mailgun.PrivateKey, conf.Mailgun.PublicKey)
}
