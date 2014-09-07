package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/majimoe/server/config"
	"github.com/majimoe/server/server"
	"gopkg.in/yaml.v1"
)

type fixture struct {
	Users []struct {
		Name     string `yaml:"name"`
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
	} `yaml:"users"`

	Domains []struct {
		Name string `yaml:"name"`
	} `yaml:"domains"`

	Records []struct {
		Name     string `yaml:"name"`
		Type     string `yaml:"type"`
		Value    string `yaml:"value"`
		TTL      int    `yaml:"ttl"`
		Priority int    `yaml:"priority"`
	} `yaml:"records"`
}

var Server *gin.Engine
var Fixture fixture

func init() {
	Server = server.Server()
	path := path.Join(config.BaseDir, "tests", "fixture.yml")
	data, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(data, &Fixture); err != nil {
		panic(err)
	}
}

type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
}

func Request(options RequestOptions) *httptest.ResponseRecorder {
	var body io.Reader

	if options.Body != nil {
		if b, err := json.Marshal(options.Body); err != nil {
			panic(err)
		} else {
			body = bytes.NewReader(b)
		}
	}

	req, err := http.NewRequest(options.Method, options.URL, body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	if options.Headers != nil {
		for key, value := range options.Headers {
			req.Header.Set(key, value)
		}
	}

	Server.ServeHTTP(w, req)

	if err != nil {
		panic(err)
	}

	return w
}

func ParseJSON(body *bytes.Buffer, data interface{}) error {
	return json.NewDecoder(body).Decode(data)
}
