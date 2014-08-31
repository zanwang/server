package util

import (
	"github.com/majimoe/server/config"
	"github.com/mrjones/oauth"
)

var Twitter *oauth.Consumer

func init() {
	conf := config.Config
	Twitter = oauth.NewConsumer(
		conf.Twitter.APIKey,
		conf.Twitter.APISecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
	)
}
