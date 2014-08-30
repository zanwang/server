package server

import (
	"net/http"

	. "github.com/onsi/gomega"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
)

func (s *TestSuite) APIv1Token() {
	s.Describe("Token", func() {
		s.APIv1TokenCreate()
	})
}

func (s *TestSuite) createToken(key string, user *models.User, body map[string]string) {
	var token models.Token
	r := s.Request("POST", "/api/v1/tokens", &requestOptions{Body: body})

	Expect(r.Code).To(Equal(http.StatusCreated))

	s.ParseJSON(r.Body, &token)
	Expect(token.Key).To(HaveLen(32))
	Expect(token.UserID).To(Equal(user.ID))

	s.Set(key, &token)
}

func (s *TestSuite) deleteToken(key string) {
	token := s.Get(key).(*models.Token)
	models.DB.Delete(token)
}

func (s *TestSuite) createToken1() {
	user := s.Get("user").(*models.User)

	s.createToken("token", user, map[string]string{
		"email":    Fixture.Users[0].Email,
		"password": Fixture.Users[0].Password,
	})
}

func (s *TestSuite) deleteToken1() {
	s.deleteToken("token")
}

func (s *TestSuite) createToken2() {
	user := s.Get("user2").(*models.User)

	s.createToken("token2", user, map[string]string{
		"email":    Fixture.Users[1].Email,
		"password": Fixture.Users[1].Password,
	})
}

func (s *TestSuite) deleteToken2() {
	s.deleteToken("token2")
}

func (s *TestSuite) APIv1TokenCreate() {
	s.Describe("Create", func() {
		s.Before(func() {
			s.createUser1()
		})

		s.It("Success", func() {
			s.createToken1()
		})

		s.It("Email required", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{Body: map[string]string{}})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("email"))
			Expect(err.Code).To(Equal(errors.Required))
			Expect(err.Message).To(Equal("Email is required"))
		})

		s.It("Password required", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{Body: map[string]string{
				"email": "abc@def.com",
			}})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("password"))
			Expect(err.Code).To(Equal(errors.Required))
			Expect(err.Message).To(Equal("Password is required"))
		})

		s.It("Email format", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{Body: map[string]string{
				"email":    "abc",
				"password": "123456",
			}})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("email"))
			Expect(err.Code).To(Equal(errors.Email))
			Expect(err.Message).To(Equal("Email is invalid"))
		})

		s.It("User does not exist", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{Body: map[string]string{
				"email":    "abc@def.com",
				"password": "123456",
			}})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("email"))
			Expect(err.Code).To(Equal(errors.UserNotExist))
			Expect(err.Message).To(Equal("User does not exist"))
		})

		s.It("Password length", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{Body: map[string]string{
				"email":    Fixture.Users[0].Email,
				"password": "123",
			}})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("password"))
			Expect(err.Code).To(Equal(errors.Length))
			Expect(err.Message).To(Equal("The length of password must be between 6-50"))
		})

		s.It("Wrong password", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{Body: map[string]string{
				"email":    Fixture.Users[0].Email,
				"password": "erqojeroqjeor",
			}})

			Expect(r.Code).To(Equal(http.StatusUnauthorized))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("password"))
			Expect(err.Code).To(Equal(errors.WrongPassword))
			Expect(err.Message).To(Equal("Password is wrong"))
		})

		s.It("Password unset", func() {
			user := s.Get("user").(*models.User)
			user.Password = ""

			if _, err := models.DB.Update(user); err != nil {
				s.Fail(err)
			}

			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{Body: map[string]string{
				"email":    Fixture.Users[0].Email,
				"password": "erqojeroqjeor",
			}})

			Expect(r.Code).To(Equal(http.StatusUnauthorized))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("password"))
			Expect(err.Code).To(Equal(errors.PasswordUnset))
			Expect(err.Message).To(Equal("Password has not been set"))
		})

		s.After(func() {
			s.deleteUser1()
			s.deleteToken1()
		})
	})
}
