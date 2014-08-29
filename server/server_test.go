package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
)

type TestSuite struct {
	*goblin.G
	server *gin.Engine
	data   map[string]interface{}
}

func (s *TestSuite) request(method, url string, data interface{}) *httptest.ResponseRecorder {
	var body io.Reader

	if data != nil {
		if b, err := json.Marshal(data); err != nil {
			panic(err)
		} else {
			body = bytes.NewReader(b)
		}
	}

	req, err := http.NewRequest(method, url, body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")
	s.server.ServeHTTP(w, req)
	Expect(err, BeNil())

	return w
}

func (s *TestSuite) parseJSON(body *bytes.Buffer, data interface{}) {
	if err := json.NewDecoder(body).Decode(data); err != nil {
		panic(err)
	}
}

type apiError struct {
	Error errors.API `json:"name"`
}

func TestServer(t *testing.T) {
	s := TestSuite{
		goblin.Goblin(t),
		Server(),
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
