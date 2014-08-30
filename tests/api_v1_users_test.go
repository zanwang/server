package server

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"

	. "github.com/onsi/gomega"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
)

func (s *TestSuite) APIv1User() {
	s.Describe("User", func() {
		s.APIv1UserCreate()
		s.APIv1UserShow()
		s.APIv1UserUpdate()
	})
}

func (s *TestSuite) createUser(key string, body map[string]string) {
	var user models.User
	r := s.Request("POST", "/api/v1/users", &requestOptions{Body: body})
	h := md5.New()

	Expect(r.Code).To(Equal(http.StatusCreated))

	s.ParseJSON(r.Body, &user)
	Expect(user.Name).To(Equal(body["name"]))
	Expect(user.Email).To(Equal(body["email"]))

	io.WriteString(h, body["email"])
	Expect(user.Avatar).To(Equal("//www.gravatar.com/avatar/" + hex.EncodeToString(h.Sum(nil))))

	s.Set(key, &user)
}

func (s *TestSuite) deleteUser(key string) {
	user := s.Get(key).(*models.User)
	models.DB.Delete(user)
	s.Del(key)
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

func (s *TestSuite) createUser2() {
	s.createUser("user2", map[string]string{
		"name":     Fixture.Users[1].Name,
		"password": Fixture.Users[1].Password,
		"email":    Fixture.Users[1].Email,
	})
}

func (s *TestSuite) deleteUser2() {
	s.deleteUser("user2")
}

func (s *TestSuite) APIv1UserCreate() {
	s.Describe("Create", func() {
		s.It("Success", func() {
			s.createUser1()
		})

		s.It("Name required", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", &requestOptions{
				Body: map[string]string{},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("name"))
			Expect(err.Code).To(Equal(errors.Required))
			Expect(err.Message).To(Equal("Name is required"))
		})

		s.It("Password required", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", &requestOptions{
				Body: map[string]string{
					"name": "John Doe",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("password"))
			Expect(err.Code).To(Equal(errors.Required))
			Expect(err.Message).To(Equal("Password is required"))
		})

		s.It("Email required", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", &requestOptions{
				Body: map[string]string{
					"name":     "John Doe",
					"password": "123456",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("email"))
			Expect(err.Code).To(Equal(errors.Required))
			Expect(err.Message).To(Equal("Email is required"))
		})

		s.It("Password length", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", &requestOptions{
				Body: map[string]string{
					"name":     "John Doe",
					"password": "123",
					"email":    "abc@def.com",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("password"))
			Expect(err.Code).To(Equal(errors.Length))
			Expect(err.Message).To(Equal("The length of password must be between 6-50"))
		})

		s.It("Email format", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", &requestOptions{
				Body: map[string]string{
					"name":     "John Doe",
					"password": "123456",
					"email":    "abc",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("email"))
			Expect(err.Code).To(Equal(errors.Email))
			Expect(err.Message).To(Equal("Email is invalid"))
		})

		s.It("Email has been taken", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", &requestOptions{
				Body: map[string]string{
					"name":     Fixture.Users[0].Name,
					"password": Fixture.Users[0].Password,
					"email":    Fixture.Users[0].Email,
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("email"))
			Expect(err.Code).To(Equal(errors.EmailUsed))
			Expect(err.Message).To(Equal("Email has been taken"))
		})

		s.After(func() {
			s.deleteUser1()
		})
	})
}

func (s *TestSuite) APIv1UserShow() {
	s.Describe("Show", func() {
		s.Before(func() {
			s.createUser1()
			s.createUser2()
			s.createToken1()
			s.createToken2()
		})

		s.It("Success (private)", func() {
			var u map[string]interface{}
			token := s.Get("token").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("GET", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusOK))

			s.ParseJSON(r.Body, &u)

			Expect(u["id"]).To(BeEquivalentTo(user.ID))
			Expect(u["name"]).To(Equal(user.Name))
			Expect(u["email"]).To(Equal(user.Email))
			Expect(u["avatar"]).To(Equal(user.Avatar))
			Expect(u["activated"]).To(Equal(user.Activated))
			Expect(u).To(HaveKey("created_at"))
			Expect(u).To(HaveKey("updated_at"))

			Expect(u).NotTo(HaveKey("password"))
			Expect(u).NotTo(HaveKey("activation_token"))
			Expect(u).NotTo(HaveKey("password_reset_token"))
			Expect(u).NotTo(HaveKey("facebook_id"))
			Expect(u).NotTo(HaveKey("twitter_id"))
			Expect(u).NotTo(HaveKey("google_id"))
			Expect(u).NotTo(HaveKey("github_id"))
		})

		s.It("Success (public)", func() {
			var u map[string]interface{}
			token := s.Get("token2").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("GET", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusOK))

			s.ParseJSON(r.Body, &u)

			Expect(u["id"]).To(BeEquivalentTo(user.ID))
			Expect(u["name"]).To(Equal(user.Name))
			Expect(u["avatar"]).To(Equal(user.Avatar))
			Expect(u).To(HaveKey("created_at"))
			Expect(u).To(HaveKey("updated_at"))

			Expect(u).NotTo(HaveKey("email"))
			Expect(u).NotTo(HaveKey("activated"))
			Expect(u).NotTo(HaveKey("password"))
			Expect(u).NotTo(HaveKey("activation_token"))
			Expect(u).NotTo(HaveKey("password_reset_token"))
			Expect(u).NotTo(HaveKey("facebook_id"))
			Expect(u).NotTo(HaveKey("twitter_id"))
			Expect(u).NotTo(HaveKey("google_id"))
			Expect(u).NotTo(HaveKey("github_id"))
		})

		s.It("Success (unauthorized)", func() {
			var u map[string]interface{}
			user := s.Get("user").(*models.User)
			r := s.Request("GET", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Body: map[string]string{},
			})

			Expect(r.Code).To(Equal(http.StatusOK))

			s.ParseJSON(r.Body, &u)

			Expect(u["id"]).To(BeEquivalentTo(user.ID))
			Expect(u["name"]).To(Equal(user.Name))
			Expect(u["avatar"]).To(Equal(user.Avatar))
			Expect(u).To(HaveKey("created_at"))
			Expect(u).To(HaveKey("updated_at"))

			Expect(u).NotTo(HaveKey("email"))
			Expect(u).NotTo(HaveKey("activated"))
			Expect(u).NotTo(HaveKey("password"))
			Expect(u).NotTo(HaveKey("activation_token"))
			Expect(u).NotTo(HaveKey("password_reset_token"))
			Expect(u).NotTo(HaveKey("facebook_id"))
			Expect(u).NotTo(HaveKey("twitter_id"))
			Expect(u).NotTo(HaveKey("google_id"))
			Expect(u).NotTo(HaveKey("github_id"))
		})

		s.It("User does not exist", func() {
			var err errors.API
			r := s.Request("GET", "/api/v1/users/46546879", &requestOptions{Body: map[string]string{}})

			Expect(r.Code).To(Equal(http.StatusNotFound))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.UserNotExist))
			Expect(err.Message).To(Equal("User does not exist"))
		})

		s.After(func() {
			s.deleteUser1()
			s.deleteUser2()
			s.deleteToken1()
			s.deleteToken2()
		})
	})
}

func (s *TestSuite) APIv1UserUpdate() {
	s.Describe("Update", func() {
		s.Before(func() {
			s.createUser1()
			s.createUser2()
			s.createToken1()
			s.createToken2()
		})

		s.It("Success", func() {
			//
		})

		s.It("Unauthorized (with token)", func() {
			//
		})

		s.It("Unauthorized (without token)", func() {
			//
		})

		s.After(func() {
			s.deleteUser1()
			s.deleteUser2()
			s.deleteToken1()
			s.deleteToken2()
		})
	})
}
