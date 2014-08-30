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
		s.APIv1UserDestroy()
	})
}

func md5str(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

func (s *TestSuite) createUser(key string, body map[string]string) {
	var user models.User
	r := s.Request("POST", "/api/v1/users", &requestOptions{Body: body})

	Expect(r.Code).To(Equal(http.StatusCreated))

	s.ParseJSON(r.Body, &user)
	Expect(user.Name).To(Equal(body["name"]))
	Expect(user.Email).To(Equal(body["email"]))
	Expect(user.Activated).To(BeFalse())
	Expect(user.Avatar).To(Equal("//www.gravatar.com/avatar/" + md5str(body["email"])))

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

		s.It("Password too short", func() {
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

		s.It("Password too long", func() {
			var err errors.API
			r := s.Request("POST", "/api/v1/users", &requestOptions{
				Body: map[string]string{
					"name":     "John Doe",
					"password": "123464646546546546546546546546546546546546546546546546546546546",
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
			var u map[string]interface{}
			token := s.Get("token").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name": "WTF",
				},
			})

			Expect(r.Code).To(Equal(http.StatusOK))

			s.ParseJSON(r.Body, &u)

			Expect(u["id"]).To(BeEquivalentTo(user.ID))
			Expect(u["name"]).To(Equal("WTF"))
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

			user.Name = u["name"].(string)
		})

		s.It("Unauthorized (with wrong token)", func() {
			var err errors.API
			token := s.Get("token2").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name": "WTF",
				},
			})

			Expect(r.Code).To(Equal(http.StatusForbidden))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.UserForbidden))
			Expect(err.Message).To(Equal("You are forbidden to access this user"))
		})

		s.It("Unauthorized (without token)", func() {
			var err errors.API
			user := s.Get("user").(*models.User)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Body: map[string]string{
					"name": "WTF",
				},
			})

			Expect(r.Code).To(Equal(http.StatusUnauthorized))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.TokenRequired))
			Expect(err.Message).To(Equal("Token is required"))
		})

		s.It("Edit email", func() {
			var u map[string]interface{}
			token := s.Get("token").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"email": "abc@def.com",
				},
			})

			Expect(r.Code).To(Equal(http.StatusOK))

			s.ParseJSON(r.Body, &u)

			Expect(u["id"]).To(BeEquivalentTo(user.ID))
			Expect(u["name"]).To(Equal("WTF"))
			Expect(u["email"]).To(Equal("abc@def.com"))
			Expect(u["avatar"]).To(Equal("//www.gravatar.com/avatar/" + md5str(u["email"].(string))))
			Expect(u["activated"]).To(BeFalse())
			Expect(u).To(HaveKey("created_at"))
			Expect(u).To(HaveKey("updated_at"))

			Expect(u).NotTo(HaveKey("password"))
			Expect(u).NotTo(HaveKey("activation_token"))
			Expect(u).NotTo(HaveKey("password_reset_token"))
			Expect(u).NotTo(HaveKey("facebook_id"))
			Expect(u).NotTo(HaveKey("twitter_id"))
			Expect(u).NotTo(HaveKey("google_id"))
			Expect(u).NotTo(HaveKey("github_id"))

			user.Email = u["email"].(string)
			user.Activated = u["activated"].(bool)
			user.Avatar = u["avatar"].(string)
		})

		s.It("Invalid email", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"email": "abc",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.Email))
			Expect(err.Message).To(Equal("Email is invalid"))
		})

		s.It("Email has been taken", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"email": Fixture.Users[1].Email,
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.EmailUsed))
			Expect(err.Message).To(Equal("Email has been taken"))
		})

		s.It("Password too short", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"old_password": Fixture.Users[0].Password,
					"password":     "123",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.Length))
			Expect(err.Message).To(Equal("The length of password must be between 6-50"))
		})

		s.It("Password too long", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"old_password": Fixture.Users[0].Password,
					"password":     "123464646546546546546546546546546546546546546546546546546546546",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.Length))
			Expect(err.Message).To(Equal("The length of password must be between 6-50"))
		})

		s.It("Wrong current password", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"old_password": "afajfjfaodjf;ad",
					"password":     "abcdef",
				},
			})

			Expect(r.Code).To(Equal(http.StatusForbidden))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.WrongPassword))
			Expect(err.Message).To(Equal("Password is wrong"))
		})

		s.It("Current password is required", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"password": "abcdef",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.Required))
			Expect(err.Message).To(Equal("Current password is required"))
		})

		s.It("Modify password", func() {
			var u map[string]interface{}
			token := s.Get("token").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"old_password": Fixture.Users[0].Password,
					"password":     "abcdef",
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

		s.It("No need of current password if current password has not been set", func() {
			// Clean current password
			user := s.Get("user").(*models.User)
			user.Password = ""
			models.DB.Update(user)

			var u map[string]interface{}
			token := s.Get("token").(*models.Token)
			r := s.Request("PUT", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"password": "abcdef",
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

		s.It("User does not exist", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			r := s.Request("PUT", "/api/v1/users/465464343545", &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

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

func (s *TestSuite) APIv1UserDestroy() {
	s.Describe("Destroy", func() {
		s.Before(func() {
			s.createUser1()
			s.createUser2()
			s.createToken1()
			s.createToken2()
		})

		s.It("Unauthorized (with wrong token)", func() {
			var err errors.API
			token := s.Get("token2").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("DELETE", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusForbidden))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.UserForbidden))
			Expect(err.Message).To(Equal("You are forbidden to access this user"))
		})

		s.It("Unauthorized (without token", func() {
			var err errors.API
			user := s.Get("user").(*models.User)
			r := s.Request("DELETE", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Body: map[string]string{},
			})

			Expect(r.Code).To(Equal(http.StatusUnauthorized))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.TokenRequired))
			Expect(err.Message).To(Equal("Token is required"))
		})

		s.It("User does not exist", func() {
			var err errors.API
			token := s.Get("token2").(*models.Token)
			r := s.Request("DELETE", "/api/v1/users/4465453135463", &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusNotFound))

			s.ParseJSON(r.Body, &err)

			Expect(err.Code).To(Equal(errors.UserNotExist))
			Expect(err.Message).To(Equal("User does not exist"))
		})

		s.It("Success", func() {
			token := s.Get("token").(*models.Token)
			user := s.Get("user").(*models.User)
			r := s.Request("DELETE", "/api/v1/users/"+strconv.FormatInt(user.ID, 10), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{},
			})

			Expect(r.Code).To(Equal(http.StatusNoContent))

			// Check whether user still exists
			if count, _ := models.DB.SelectInt("SELECT * FROM users WHERE id=?", user.ID); count > 0 {
				s.Fail("User still exists")
			}
		})

		s.After(func() {
			s.deleteUser1()
			s.deleteUser2()
			s.deleteToken1()
			s.deleteToken2()
		})
	})
}
