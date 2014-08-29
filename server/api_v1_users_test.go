package server

import (
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

func (s *TestSuite) createUser() {
	var user models.User
	r := s.request("POST", "/api/v1/users", map[string]interface{}{
		"name":     "John Doe",
		"password": "123456",
		"email":    "abc@def.com",
	})

	s.parseJSON(r.Body, &user)
	Expect(r.Code, http.StatusCreated)
	Expect(user.Name, "John Doe")
	Expect(user.Email, "abc@def.com")
	Expect(user.Avatar, "//www.gravatar.com/avatar/b188d046267bb5cddbc457580551297d")

	s.data["user"] = &user
}

func (s *TestSuite) deleteUser() {
	user := s.data["user"].(*models.User)
	models.DB.Delete(user)
}

func (s *TestSuite) APIv1UserCreate() {
	s.Describe("Create", func() {
		s.It("Success", func() {
			s.createUser()
		})

		s.It("Name required", func() {
			var err apiError
			r := s.request("POST", "/api/v1/users", nil)

			s.parseJSON(r.Body, &err)
			Expect(r.Code, http.StatusBadRequest)
			Expect(err.Error.Field, "name")
			Expect(err.Error.Code, errors.Required)
			Expect(err.Error.Message, "Name is required")
		})

		s.It("Password required", func() {
			var err apiError
			r := s.request("POST", "/api/v1/users", map[string]interface{}{
				"name": "John Doe",
			})

			s.parseJSON(r.Body, &err)
			Expect(r.Code, http.StatusBadRequest)
			Expect(err.Error.Field, "password")
			Expect(err.Error.Code, errors.Required)
			Expect(err.Error.Message, "Password is required")
		})

		s.It("Email required", func() {
			var err apiError
			r := s.request("POST", "/api/v1/users", map[string]interface{}{
				"name":     "John Doe",
				"password": "123456",
			})

			s.parseJSON(r.Body, &err)
			Expect(r.Code, http.StatusBadRequest)
			Expect(err.Error.Field, "email")
			Expect(err.Error.Code, errors.Required)
			Expect(err.Error.Message, "Email is required")
		})

		s.It("Password length", func() {
			var err apiError
			r := s.request("POST", "/api/v1/users", map[string]interface{}{
				"name":     "John Doe",
				"password": "123",
				"email":    "abc@def.com",
			})

			s.parseJSON(r.Body, &err)
			Expect(r.Code, http.StatusBadRequest)
			Expect(err.Error.Field, "password")
			Expect(err.Error.Code, errors.Length)
			Expect(err.Error.Message, "The length of password must be between 6-50")
		})

		s.It("Email format", func() {
			var err apiError
			r := s.request("POST", "/api/v1/users", map[string]interface{}{
				"name":     "John Doe",
				"password": "123456",
				"email":    "abc",
			})

			s.parseJSON(r.Body, &err)
			Expect(r.Code, http.StatusBadRequest)
			Expect(err.Error.Field, "email")
			Expect(err.Error.Code, errors.Email)
			Expect(err.Error.Message, "Email is invalid")
		})

		s.After(func() {
			s.deleteUser()
		})
	})
}
