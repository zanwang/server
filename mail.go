package main

import (
	"github.com/mailgun/mailgun-go"
	"github.com/tommy351/maji.moe/config"
)

func mail(config *config.Config) mailgun.Mailgun {
	mg := mailgun.NewMailgun(config.Mailgun.Domain, config.Mailgun.PrivateKey, config.Mailgun.PublicKey)
	return mg
}
