package server

import (
	"net/http"

	. "github.com/onsi/gomega"
	"github.com/tommy351/maji.moe/models"
)

func (s *TestSuite) APIv1Token() {
	s.Describe("Token", func() {
		s.APIv1TokenCreate()
	})
}

func (s *TestSuite) createToken() {
	var token models.Token
	user := s.data["user"].(*models.User)
	r := s.request("POST", "/api/v1/tokens", map[string]interface{}{
		"email":    "abc@def.com",
		"password": "123456",
	})

	s.parseJSON(r.Body, &token)
	Expect(r.Code, http.StatusCreated)
	Expect(token.Key, HaveLen(32))
	Expect(token.UserID, user.ID)

	s.data["token"] = &token
}

func (s *TestSuite) deleteToken() {
	token := s.data["token"].(*models.Token)
	models.DB.Delete(token)
}

func (s *TestSuite) APIv1TokenCreate() {
	s.Describe("Create", func() {
		s.Before(func() {
			s.createUser()
		})

		s.It("Success", func() {
			s.createToken()
		})

		s.It("Email required", func() {
			//
		})

		s.It("Password required", func() {
			//
		})

		s.It("Email format", func() {
			//
		})

		s.It("Password length", func() {
			//
		})

		s.It("Wrong password", func() {
			//
		})

		s.It("Password unset", func() {
			//
		})

		s.After(func() {
			s.deleteUser()
			s.deleteToken()
		})
	})
}
