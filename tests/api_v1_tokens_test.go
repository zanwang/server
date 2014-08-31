package tests

import (
	"net/http"
	"time"

	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/models"
	. "github.com/onsi/gomega"
)

func (s *TestSuite) APIv1Token() {
	s.Describe("Token", func() {
		s.APIv1TokenCreate()
		s.APIv1TokenUpdate()
		s.APIv1TokenDestroy()
	})
}

func (s *TestSuite) createToken(key string, user *models.User, body map[string]string) {
	var token models.Token
	r := s.Request("POST", "/api/v1/tokens", &requestOptions{Body: body})

	Expect(r.Code).To(Equal(http.StatusCreated))
	Expect(r.Header().Get("Pragma")).To(Equal("no-cache"))
	Expect(r.Header().Get("Cache-Control")).To(Equal("no-cache, no-store, must-revalidate"))
	Expect(r.Header().Get("Expires")).To(Equal("0"))

	s.ParseJSON(r.Body, &token)
	Expect(token.Key).To(HaveLen(32))
	Expect(token.UserID).To(Equal(user.ID))
	Expect(token.UpdatedAt.Add(time.Hour * 24 * 7)).To(Equal(token.ExpiredAt))

	s.Set(key, &token)
}

func (s *TestSuite) deleteToken(key string) {
	token := s.Get(key).(*models.Token)
	models.DB.Delete(token)
	s.Del(key)
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
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{
				Body: map[string]string{},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("email"))
			Expect(err.Code).To(Equal(errors.Required))
			Expect(err.Message).To(Equal("Email is required"))
		})

		s.It("Password required", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{
				Body: map[string]string{
					"email": "abc@def.com",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("password"))
			Expect(err.Code).To(Equal(errors.Required))
			Expect(err.Message).To(Equal("Password is required"))
		})

		s.It("Invalid email", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{
				Body: map[string]string{
					"email":    "abc",
					"password": "123456",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("email"))
			Expect(err.Code).To(Equal(errors.Email))
			Expect(err.Message).To(Equal("Email is invalid"))
		})

		s.It("User does not exist", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{
				Body: map[string]string{
					"email":    "abc@def.com",
					"password": "123456",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("email"))
			Expect(err.Code).To(Equal(errors.UserNotExist))
			Expect(err.Message).To(Equal("User does not exist"))
		})

		s.It("Password too short", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{
				Body: map[string]string{
					"email":    Fixture.Users[0].Email,
					"password": "123",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("password"))
			Expect(err.Code).To(Equal(errors.Length))
			Expect(err.Message).To(Equal("The length of password must be between 6-50"))
		})

		s.It("Password too long", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{
				Body: map[string]string{
					"email":    Fixture.Users[0].Email,
					"password": "123464646546546546546546546546546546546546546546546546546546546",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("password"))
			Expect(err.Code).To(Equal(errors.Length))
			Expect(err.Message).To(Equal("The length of password must be between 6-50"))
		})

		s.It("Wrong password", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{
				Body: map[string]string{
					"email":    Fixture.Users[0].Email,
					"password": "erqojeroqjeor",
				},
			})

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
			r := s.Request("POST", "/api/v1/tokens", &requestOptions{
				Body: map[string]string{
					"email":    Fixture.Users[0].Email,
					"password": "erqojeroqjeor",
				},
			})

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

func (s *TestSuite) APIv1TokenUpdate() {
	s.Describe("Update", func() {
		s.Before(func() {
			s.createUser1()
			s.createToken1()
		})

		s.It("Success", func() {
			var t map[string]interface{}
			token := s.Get("token").(*models.Token)
			r := s.Request("PUT", "/api/v1/tokens", &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusOK))
			Expect(r.Header().Get("Pragma")).To(Equal("no-cache"))
			Expect(r.Header().Get("Cache-Control")).To(Equal("no-cache, no-store, must-revalidate"))
			Expect(r.Header().Get("Expires")).To(Equal("0"))

			s.ParseJSON(r.Body, &t)

			Expect(t["key"]).To(Equal(token.Key))
			Expect(t["user_id"]).To(BeEquivalentTo(token.UserID))
		})

		s.It("Unauthorized (without token)", func() {
			var err errors.API
			r := s.Request("PUT", "/api/v1/tokens", &requestOptions{})

			Expect(r.Code).To(Equal(http.StatusUnauthorized))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.TokenRequired))
			Expect(err.Message).To(Equal("Token is required"))
		})

		s.After(func() {
			s.deleteUser1()
			s.deleteToken1()
		})
	})
}

func (s *TestSuite) APIv1TokenDestroy() {
	s.Describe("Destroy", func() {
		s.Before(func() {
			s.createUser1()
			s.createToken1()
		})

		s.It("Success", func() {
			token := s.Get("token").(*models.Token)
			r := s.Request("DELETE", "/api/v1/tokens", &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusNoContent))

			// Check whether token still exists
			if count, _ := models.DB.SelectInt("SELECT count(*) FROM tokens WHERE key=?", token.Key); count > 0 {
				s.Fail("Token still exists")
			}
		})

		s.It("Unauthorized (without token)", func() {
			var err errors.API
			r := s.Request("DELETE", "/api/v1/tokens", &requestOptions{})

			Expect(r.Code).To(Equal(http.StatusUnauthorized))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.TokenRequired))
			Expect(err.Message).To(Equal("Token is required"))
		})

		s.After(func() {
			s.deleteUser1()
			s.deleteToken1()
		})
	})
}
