package oauth

import (
	"github.com/huandu/facebook"
	"github.com/tommy351/maji.moe/config"
)

func LoadFacebook(config *config.Config) *facebook.App {
	return facebook.New(config.Facebook.AppID, config.Facebook.AppSecret)
}
