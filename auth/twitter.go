package auth

import (
	"github.com/mrjones/oauth"
	"github.com/tommy351/maji.moe/config"
)

type TwitterConsumer struct {
	*oauth.Consumer
}

func LoadTwitter(config *config.Config) *TwitterConsumer {
	c := oauth.NewConsumer(
		config.Twitter.APIKey,
		config.Twitter.APISecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
	)

	return &TwitterConsumer{c}
}
