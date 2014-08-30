package server

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"

	. "github.com/onsi/gomega"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
)

func (s *TestSuite) APIv1User() {
	s.Describe("User", func() {
		s.APIv1UserCreate()
	})
}

func (s *TestSuite) createUser(key string, body map[string]string) {
	var user models.User
	r := s.Request("POST", "/api/v1/users", &requestOptions{Body: body})
	h := md5.New()

	Expect(r.Code, http.StatusCreated)

	s.ParseJSON(r.Body, &user)
	Expect(user.Name, body["name"])
	Expect(user.Email, body["email"])

	io.WriteString(h, body["email"])
	Expect(user.Avatar, "//www.gravatar.com/avatar/"+hex.EncodeToString(h.Sum(nil)))

	s.Set(key, &user)
}

func (s *TestSuite) deleteUser(key string) {
	user := s.Get(key).(*models.User)
	models.DB.Delete(user)
}

func (s *TestSuite) createUser1() {
	s.createUser("user", map[string]string{
		"name":     Fixture.Users[0].Name,
		"password": Fixture.Users[0].Password,
		"email":    Fixture.Users[0].Email,
	})
}

func (s *TestSuite) deleteUser1() {
	s.deleteUser("user")
}

func (s *TestSuite) APIv1UserCreate() {
	s.Describe("Create", func() {
		s.It("Success", func() {
			s.createUser1()
		})

		s.It("Name required", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", nil)

			s.ParseJSON(r.Body, &err)
			Expect(r.Code, http.StatusBadRequest)
			Expect(err.Field, "name")
			Expect(err.Code, errors.Required)
			Expect(err.Message, "Name is required")
		})

		s.It("Password required", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", &requestOptions{Body: map[string]string{
				"name": "John Doe",
			}})

			s.ParseJSON(r.Body, &err)
			Expect(r.Code, http.StatusBadRequest)
			Expect(err.Field, "password")
			Expect(err.Code, errors.Required)
			Expect(err.Message, "Password is required")
		})

		s.It("Email required", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", &requestOptions{Body: map[string]string{
				"name":     "John Doe",
				"password": "123456",
			}})

			s.ParseJSON(r.Body, &err)
			Expect(r.Code, http.StatusBadRequest)
			Expect(err.Field, "email")
			Expect(err.Code, errors.Required)
			Expect(err.Message, "Email is required")
		})

		s.It("Password length", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", &requestOptions{Body: map[string]string{
				"name":     "John Doe",
				"password": "123",
				"email":    "abc@def.com",
			}})

			s.ParseJSON(r.Body, &err)
			Expect(r.Code, http.StatusBadRequest)
			Expect(err.Field, "password")
			Expect(err.Code, errors.Length)
			Expect(err.Message, "The length of password must be between 6-50")
		})

		s.It("Email format", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", &requestOptions{Body: map[string]string{
				"name":     "John Doe",
				"password": "123456",
				"email":    "abc",
			}})

			s.ParseJSON(r.Body, &err)
			Expect(r.Code, http.StatusBadRequest)
			Expect(err.Field, "email")
			Expect(err.Code, errors.Email)
			Expect(err.Message, "Email is invalid")
		})

		s.After(func() {
			s.deleteUser1()
		})
	})
}
