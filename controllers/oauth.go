package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/martini-contrib/sessions"
	"github.com/mrjones/oauth"
	"github.com/tommy351/maji.moe/auth"
)

func OAuthTwitterLogin(twitter *auth.TwitterConsumer, r *http.Request, session sessions.Session, w http.ResponseWriter) {
	callback := fmt.Sprintf("http://%s/oauth/twitter/callback", r.Host)

	token, url, err := twitter.GetRequestTokenAndUrl(callback)

	if err != nil {
		panic(err)
	}

	session.Set("oauth_secret", token.Secret)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

type twitterCredential struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func OAuthTwitterCallback(twitter *auth.TwitterConsumer, r *http.Request, session sessions.Session, w http.ResponseWriter) {
	query := r.URL.Query()
	secret := session.Get("oauth_secret")

	if secret == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := twitter.AuthorizeToken(&oauth.RequestToken{
		Token:  query.Get("oauth_token"),
		Secret: secret.(string),
	}, query.Get("oauth_verifier"))

	if err != nil {
		panic(err)
	}

	res, err := twitter.Get("https://api.twitter.com/1.1/account/verify_credentials.json", nil, token)

	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(res.Body)
}
