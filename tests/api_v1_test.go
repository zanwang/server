package server

import (
	"net/http"

	. "github.com/onsi/gomega"
)

func (s *TestSuite) APIv1() {
	s.Describe("API v1", func() {
		s.APIv1Entry()
		s.APIv1User()
		s.APIv1Token()
	})
}

func (s *TestSuite) APIv1Entry() {
	s.It("Entry", func() {
		var data map[string]interface{}
		r := s.Request("GET", "/api/v1", nil)

		s.ParseJSON(r.Body, &data)
		Expect(r.Code, http.StatusOK)
		Expect(data, map[string]interface{}{
			"tokens":  "/api/v1/tokens",
			"users":   "/api/v1/users",
			"domains": "/api/v1/domains",
			"records": "/api/v1/records",
			"emails":  "/api/v1/emails",
		})
	})
}
