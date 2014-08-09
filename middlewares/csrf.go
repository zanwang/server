package middlewares

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"github.com/tommy351/maji.moe/config"
)

// CSRFToken stores CSRF token
type CSRFToken struct {
	Hash    string
	Secret  string
	Token   string
	session sessions.Session
}

// GetToken returns a valid CSRF token
func (c *CSRFToken) GetToken() string {
	if c.Token != "" {
		return c.Token
	}

	c.Token = csrfTokenize(c.Hash, c.Secret)
	c.session.Set("csrfSecret", c.Hash)

	return c.Token
}

func csrfTokenize(hash string, secret string) string {
	h := sha1.New()

	io.WriteString(h, hash)
	io.WriteString(h, secret)

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// CSRFMiddleware is a middleware that checks CSRF token for each requests
func CSRFMiddleware(c martini.Context, s sessions.Session, req *http.Request, config *config.Config, res http.ResponseWriter) {
	token := CSRFToken{
		Hash:    uniuri.New(),
		Secret:  config.Server.Secret,
		session: s,
	}

	c.Map(&token)

	switch req.Method {
	case "GET", "HEAD", "OPTIONS":
		return
	}

	hash := s.Get("csrfSecret").(string)

	if csrfTokenize(hash, config.Server.Secret) != req.FormValue("_csrf") {
		http.Error(res, "Invalid csrf token.", http.StatusBadRequest)
	}
}
