package server

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
	"github.com/tommy351/maji.moe/config"
	"github.com/tommy351/maji.moe/models"
	"github.com/tommy351/maji.moe/server"
	"gopkg.in/yaml.v1"
)

type fixture struct {
	Users []struct {
		Name     string `yaml:"name"`
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
	} `yaml:"users"`
}

var Fixture fixture

func init() {
	path := path.Join(config.BaseDir, "tests", "fixture.yml")
	data, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(data, &Fixture); err != nil {
		panic(err)
	}
}

type TestSuite struct {
	*goblin.G
	server *gin.Engine
	data   map[string]interface{}
}

type requestOptions struct {
	Headers map[string]string
	Body    interface{}
}

func (s *TestSuite) Request(method, url string, options *requestOptions) *httptest.ResponseRecorder {
	var body io.Reader

	if options != nil {
		if options.Body != nil {
			if b, err := json.Marshal(options.Body); err != nil {
				panic(err)
			} else {
				body = bytes.NewReader(b)
			}
		}
	}

	req, err := http.NewRequest(method, url, body)
	w := httptest.NewRecorder()

	if options != nil && options.Headers != nil {
		for key, value := range options.Headers {
			req.Header.Set(key, value)
		}
	}

	req.Header.Set("Content-Type", "application/json")
	s.server.ServeHTTP(w, req)
	Expect(err).To(BeNil())

	return w
}

func (s *TestSuite) ParseJSON(body *bytes.Buffer, data interface{}) {
	if err := json.NewDecoder(body).Decode(data); err != nil {
		s.Fail(err)
	}
}

func (s *TestSuite) ParseBody(body *bytes.Buffer) []byte {
	if data, err := ioutil.ReadAll(body); err == nil {
		return data
	} else {
		s.Fail(err)
	}

	return nil
}

func (s *TestSuite) Get(key string) interface{} {
	return s.data[key]
}

func (s *TestSuite) Set(key string, data interface{}) {
	s.data[key] = data
}

func (s *TestSuite) Del(key string) {
	s.data[key] = nil
}

func TestServer(t *testing.T) {
	s := TestSuite{
		goblin.Goblin(t),

		server.Server(),
		map[string]interface{}{},
	}

	RegisterFailHandler(func(m string, _ ...int) {
		s.Fail(m)
	})

	s.Describe("Server test", func() {
		s.APIv1()

		s.After(func() {
			models.DB.DropTables()
		})
	})
}
