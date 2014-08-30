package tests

import (
	"net/http"
	"strconv"
	"time"

	. "github.com/onsi/gomega"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
)

func (s *TestSuite) APIv1Domain() {
	s.Describe("Domain", func() {
		s.APIv1DomainCreate()
		s.APIv1DomainList()
		s.APIv1DomainShow()
	})
}

func domainCreateURL(id int64) string {
	return "/api/v1/users/" + strconv.FormatInt(id, 10) + "/domains"
}

func (s *TestSuite) createDomain(key string, token *models.Token, body map[string]string) {
	var domain models.Domain
	r := s.Request("POST", domainCreateURL(token.UserID), &requestOptions{
		Body: body,
		Headers: map[string]string{
			"Authorization": "token " + token.Key,
		},
	})

	Expect(r.Code).To(Equal(http.StatusCreated))

	s.ParseJSON(r.Body, &domain)

	Expect(domain.Name).To(Equal(body["name"]))
	Expect(domain.UserID).To(Equal(token.UserID))
	Expect(domain.CreatedAt.Add(time.Hour * 24 * 365)).To(Equal(domain.ExpiredAt))

	s.Set(key, &domain)
}

func (s *TestSuite) deleteDomain(key string) {
	domain := s.Get(key).(*models.Domain)
	models.DB.Delete(domain)
	s.Del(key)
}

func (s *TestSuite) createDomain1() {
	token := s.Get("token").(*models.Token)
	s.setUserActivated("user", true)
	s.createDomain("domain", token, map[string]string{
		"name": Fixture.Domains[0].Name,
	})
}

func (s *TestSuite) deleteDomain1() {
	s.deleteDomain("domain")
}

func (s *TestSuite) createDomain2() {
	token := s.Get("token2").(*models.Token)
	s.setUserActivated("user2", true)
	s.createDomain("domain2", token, map[string]string{
		"name": Fixture.Domains[1].Name,
	})
}

func (s *TestSuite) deleteDomain2() {
	s.deleteDomain("domain2")
}

func (s *TestSuite) APIv1DomainCreate() {
	s.Describe("Create", func() {
		s.Before(func() {
			s.createUser1()
			s.createUser2()
			s.createToken1()
			s.createToken2()
			s.createDomain1()
		})

		s.It("User has not been activated", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			s.setUserActivated("user", false)

			r := s.Request("POST", domainCreateURL(token.UserID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name": "abc",
				},
			})

			Expect(r.Code).To(Equal(http.StatusForbidden))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.UserNotActivated))
			Expect(err.Message).To(Equal("User has not been activated"))
		})

		s.It("Success", func() {
			s.createDomain2()
		})

		s.It("Name required", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			s.setUserActivated("user", true)

			r := s.Request("POST", domainCreateURL(token.UserID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.Required))
			Expect(err.Message).To(Equal("Name is required"))
		})

		s.It("Domain name started with number", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			s.setUserActivated("user", true)

			r := s.Request("POST", domainCreateURL(token.UserID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name": "1a2b",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("name"))
			Expect(err.Code).To(Equal(errors.DomainName))
			Expect(err.Message).To(Equal("Domain name is invalid"))
		})

		s.It("Domain name with special characters", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			s.setUserActivated("user", true)

			r := s.Request("POST", domainCreateURL(token.UserID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name": "中文",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("name"))
			Expect(err.Code).To(Equal(errors.DomainName))
			Expect(err.Message).To(Equal("Domain name is invalid"))
		})

		s.It("Domain name too long", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			s.setUserActivated("user", true)

			r := s.Request("POST", domainCreateURL(token.UserID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name": "erfwoerjweijrliejrwejrliwejrwjerliwwjeroiljweloirjweolirireorweorjweorwoerjwoerwoeirlsj",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("name"))
			Expect(err.Code).To(Equal(errors.MaxLength))
			Expect(err.Message).To(Equal("Maximum length of name is 63"))
		})

		s.It("Forbidden (with wrong token)", func() {
			var err errors.API
			token := s.Get("token2").(*models.Token)
			user := s.Get("user").(*models.User)
			s.setUserActivated("user", true)

			r := s.Request("POST", domainCreateURL(user.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name": "abc",
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
			s.setUserActivated("user", true)

			r := s.Request("POST", domainCreateURL(user.ID), &requestOptions{
				Body: map[string]string{
					"name": "abc",
				},
			})

			Expect(r.Code).To(Equal(http.StatusUnauthorized))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.TokenRequired))
			Expect(err.Message).To(Equal("Token is required"))
		})

		s.It("Domain name has been taken", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			s.setUserActivated("user", true)

			r := s.Request("POST", domainCreateURL(token.UserID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name": Fixture.Domains[0].Name,
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("name"))
			Expect(err.Code).To(Equal(errors.DomainUsed))
			Expect(err.Message).To(Equal("Domain name has been taken"))
		})

		s.It("Domain name has been reserved", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			s.setUserActivated("user", true)

			r := s.Request("POST", domainCreateURL(token.UserID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name": "www",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("name"))
			Expect(err.Code).To(Equal(errors.DomainReserved))
			Expect(err.Message).To(Equal("Domain name has been reserved"))
		})

		s.After(func() {
			s.deleteUser1()
			s.deleteUser2()
			s.deleteToken1()
			s.deleteToken2()
			s.deleteDomain1()
			s.deleteDomain2()
		})
	})
}

func (s *TestSuite) APIv1DomainList() {
	s.Describe("List", func() {
		s.Before(func() {
			s.createUser1()
			s.createUser2()
			s.createToken1()
			s.createToken2()
			s.createDomain1()
			s.createDomain2()
		})

		s.It("Success (owner)", func() {
			var domains []models.Domain
			domain := s.Get("domain").(*models.Domain)
			token := s.Get("token").(*models.Token)
			r := s.Request("GET", domainCreateURL(token.UserID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusOK))

			s.ParseJSON(r.Body, &domains)
			d := domains[0]
			Expect(d.ID).To(Equal(domain.ID))
			Expect(d.Name).To(Equal(domain.Name))
			Expect(d.UserID).To(Equal(domain.UserID))
			Expect(d.CreatedAt.Add(time.Hour * 24 * 365)).To(Equal(d.ExpiredAt))
		})

		s.It("Success (others)", func() {
			var domains []models.Domain
			domain := s.Get("domain2").(*models.Domain)
			token := s.Get("token").(*models.Token)
			user := s.Get("user2").(*models.User)
			r := s.Request("GET", domainCreateURL(user.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusOK))

			s.ParseJSON(r.Body, &domains)
			d := domains[0]
			Expect(d.ID).To(Equal(domain.ID))
			Expect(d.Name).To(Equal(domain.Name))
			Expect(d.UserID).To(Equal(domain.UserID))
			Expect(d.CreatedAt.Add(time.Hour * 24 * 365)).To(Equal(d.ExpiredAt))
		})

		s.It("Success (guest)", func() {
			var domains []models.Domain
			domain := s.Get("domain").(*models.Domain)
			token := s.Get("token").(*models.Token)

			r := s.Request("GET", domainCreateURL(token.UserID), nil)

			Expect(r.Code).To(Equal(http.StatusOK))

			s.ParseJSON(r.Body, &domains)
			d := domains[0]
			Expect(d.ID).To(Equal(domain.ID))
			Expect(d.Name).To(Equal(domain.Name))
			Expect(d.UserID).To(Equal(domain.UserID))
			Expect(d.CreatedAt.Add(time.Hour * 24 * 365)).To(Equal(d.ExpiredAt))
		})

		s.It("User does not exist", func() {
			var err errors.API
			r := s.Request("GET", domainCreateURL(9999999999), nil)

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
			s.deleteDomain1()
			s.deleteDomain2()
		})
	})
}

func (s *TestSuite) APIv1DomainShow() {
	s.Describe("Show", func() {
		s.Before(func() {
			s.createUser1()
			s.createToken1()
			s.createDomain1()
		})

		s.It("Success", func() {
			var d map[string]interface{}
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("GET", "/api/v1/domains/"+strconv.FormatInt(domain.ID, 10), nil)

			Expect(r.Code).To(Equal(http.StatusOK))

			s.ParseJSON(r.Body, &d)
			Expect(d["id"]).To(BeEquivalentTo(domain.ID))
			Expect(d["name"]).To(Equal(domain.Name))
			Expect(d["user_id"]).To(BeEquivalentTo(domain.UserID))
			Expect(d).To(HaveKey("created_at"))
			Expect(d).To(HaveKey("updated_at"))
			Expect(d).To(HaveKey("expired_at"))
		})

		s.It("Domain does not exist", func() {
			var err errors.API
			r := s.Request("GET", "/api/v1/domains/9999999999", nil)

			Expect(r.Code).To(Equal(http.StatusNotFound))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.DomainNotExist))
			Expect(err.Message).To(Equal("Domain does not exist"))
		})

		s.After(func() {
			s.deleteUser1()
			s.deleteToken1()
			s.deleteDomain1()
		})
	})
}
