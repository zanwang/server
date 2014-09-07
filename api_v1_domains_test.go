package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/models"
	. "github.com/smartystreets/goconvey/convey"
)

func createDomain(token *models.Token, body map[string]interface{}) (*httptest.ResponseRecorder, *models.Domain) {
	var domain models.Domain
	r := Request(RequestOptions{
		Method: "POST",
		URL:    "/api/v1/users/" + strconv.FormatInt(token.UserId, 10) + "/domains",
		Headers: map[string]string{
			"Authorization": "token " + token.Key,
		},
		Body: body,
	})

	if err := ParseJSON(r.Body, &domain); err != nil {
		panic(err)
	}

	return r, &domain
}

func createDomain1(token *models.Token) (*httptest.ResponseRecorder, *models.Domain) {
	return createDomain(token, map[string]interface{}{
		"name": Fixture.Domains[0].Name,
	})
}

func createDomain2(token *models.Token) (*httptest.ResponseRecorder, *models.Domain) {
	return createDomain(token, map[string]interface{}{
		"name": Fixture.Domains[1].Name,
	})
}

func deleteDomain(domain *models.Domain) {
	if err := models.DB.Delete(domain).Error; err != nil {
		panic(err)
	}
}

func domainCreateURL(user *models.User) string {
	return "/api/v1/users/" + strconv.FormatInt(user.Id, 10) + "/domains"
}

func TestAPIv1DomainCreate(t *testing.T) {
	Convey("API v1 - Domain create", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()

		Convey("Success", func() {
			activateUser(u1)

			r, domain := createDomain1(t1)
			defer deleteDomain(domain)

			So(r.Code, ShouldEqual, http.StatusCreated)
			So(domain.Name, ShouldEqual, Fixture.Domains[0].Name)
			So(domain.CreatedAt, ShouldNotBeNil)
			So(domain.UpdatedAt, ShouldNotBeNil)
			So(domain.ExpiredAt, ShouldResemble, domain.CreatedAt.AddDate(1, 0, 0))
		})

		Convey("Name required", func() {
			activateUser(u1)

			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    domainCreateURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "name")
			So(err.Code, ShouldEqual, errors.Required)
			So(err.Message, ShouldEqual, "Name is required")
		})

		Convey("Domain name with special characters", func() {
			activateUser(u1)

			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    domainCreateURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"name": "中文",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "name")
			So(err.Code, ShouldEqual, errors.DomainName)
			So(err.Message, ShouldEqual, "Only numbers and characters are allowed in domain name")
		})

		Convey("Domain name too long", func() {
			activateUser(u1)

			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    domainCreateURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"name": "erfwoerjweijrliejrwejrliwejrwjerliwwjeroiljweloirjweolirireorweorjweorwoerjwoerwoeirlsj",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "name")
			So(err.Code, ShouldEqual, errors.MaxLength)
			So(err.Message, ShouldEqual, "Maximum length of name is 63")
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    domainCreateURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t2.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusForbidden)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.UserForbidden)
			So(err.Message, ShouldEqual, "You are forbidden to access this user")
		})

		Convey("Unauthorized", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    domainCreateURL(u1),
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})

		Convey("Domain name has been taken", func() {
			activateUser(u1)
			_, domain := createDomain1(t1)
			defer deleteDomain(domain)

			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    domainCreateURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"name": Fixture.Domains[0].Name,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "name")
			So(err.Code, ShouldEqual, errors.DomainUsed)
			So(err.Message, ShouldEqual, "Domain name has been taken")
		})

		Convey("Domain name has been reserved", func() {
			activateUser(u1)

			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    domainCreateURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"name": "www",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "name")
			So(err.Code, ShouldEqual, errors.DomainReserved)
			So(err.Message, ShouldEqual, "Domain name has been reserved")
		})

		Convey("User does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/users/0/domains",
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.UserNotExist)
			So(err.Message, ShouldEqual, "User does not exist")
		})

		Reset(func() {
			deleteUser(u1)
			deleteUser(u2)
			deleteToken(t1)
			deleteToken(t2)
		})
	})
}

func TestAPIv1DomainList(t *testing.T) {
	Convey("API v1 - Domain list", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()
		activateUser(u1)
		activateUser(u2)
		_, d1 := createDomain1(t1)

		Convey("Success", func() {
			var domains []models.Domain
			r := Request(RequestOptions{
				Method: "GET",
				URL:    domainCreateURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &domains)
			So(len(domains), ShouldEqual, 1)
			So(domains[0], ShouldResemble, *d1)
		})

		Convey("User does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "GET",
				URL:    "/api/v1/users/0/domains",
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.UserNotExist)
			So(err.Message, ShouldEqual, "User does not exist")
		})

		Reset(func() {
			deleteUser(u1)
			deleteUser(u2)
			deleteToken(t1)
			deleteToken(t2)
			deleteDomain(d1)
		})
	})
}

func domainURL(domain *models.Domain) string {
	return "/api/v1/domains/" + strconv.FormatInt(domain.Id, 10)
}

func TestAPIv1DomainShow(t *testing.T) {
	Convey("API v1 - Domain show", t, func() {
		_, u1 := createUser1()
		_, t1 := createToken1()
		activateUser(u1)
		_, d1 := createDomain1(t1)

		Convey("Success", func() {
			var domain models.Domain
			r := Request(RequestOptions{
				Method: "GET",
				URL:    domainURL(d1),
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &domain)
			So(domain, ShouldResemble, *d1)
		})

		Convey("Domain does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "GET",
				URL:    "/api/v1/domains/0",
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.DomainNotExist)
			So(err.Message, ShouldEqual, "Domain does not exist")
		})

		Reset(func() {
			deleteUser(u1)
			deleteToken(t1)
			deleteDomain(d1)
		})
	})
}

func TestAPIv1DomainUpdate(t *testing.T) {
	Convey("API v1 - Domain update", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()
		activateUser(u1)
		activateUser(u2)
		_, d1 := createDomain1(t1)
		_, d2 := createDomain2(t2)

		Convey("Success", func() {
			time.Sleep(time.Second)

			var domain models.Domain
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    domainURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"name": "foo",
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &domain)
			So(domain.Id, ShouldEqual, d1.Id)
			So(domain.Name, ShouldEqual, "foo")
			So(domain.CreatedAt, ShouldResemble, d1.CreatedAt)
			So(domain.UpdatedAt, ShouldHappenAfter, d1.UpdatedAt)
			So(domain.ExpiredAt, ShouldResemble, d1.ExpiredAt)
			So(domain.UserId, ShouldEqual, d1.UserId)
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    domainURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t2.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusForbidden)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.DomainForbidden)
			So(err.Message, ShouldEqual, "You are forbidden to access this domain")
		})

		Convey("Unauthorized", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    domainURL(d1),
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})

		Convey("Domain name with special characters", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    domainURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]string{
					"name": "中文",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.DomainName)
			So(err.Message, ShouldEqual, "Only numbers and characters are allowed in domain name")
		})

		Convey("Domain name too long", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    domainURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]string{
					"name": "erfwoerjweijrliejrwejrliwejrwjerliwwjeroiljweloirjweolirireorweorjweorwoerjwoerwoeirlsj",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.MaxLength)
			So(err.Message, ShouldEqual, "Maximum length of name is 63")
		})

		Convey("Domain name has been taken", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    domainURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]string{
					"name": d2.Name,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.DomainUsed)
			So(err.Message, ShouldEqual, "Domain name has been taken")
		})

		Convey("Domain name has been reserved", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    domainURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]string{
					"name": "www",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.DomainReserved)
			So(err.Message, ShouldEqual, "Domain name has been reserved")
		})

		Convey("Domain name does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    "/api/v1/domains/0",
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.DomainNotExist)
			So(err.Message, ShouldEqual, "Domain does not exist")
		})

		Reset(func() {
			deleteUser(u1)
			deleteUser(u2)
			deleteToken(t1)
			deleteToken(t2)
			deleteDomain(d1)
			deleteDomain(d2)
		})
	})
}

func TestAPIv1DomainDestroy(t *testing.T) {
	Convey("API v1 - Domain destroy", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()
		activateUser(u1)
		_, d1 := createDomain1(t1)

		Convey("Success", func() {
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    domainURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNoContent)

			// Confirm domain has been deleted
			var count int

			if err := models.DB.Table("domains").Where("id = ?", d1.Id).Count(&count).Error; err != nil {
				panic(err)
			}

			So(count, ShouldEqual, 0)
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    domainURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t2.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusForbidden)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.DomainForbidden)
			So(err.Message, ShouldEqual, "You are forbidden to access this domain")
		})

		Convey("Unauthorized", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    domainURL(d1),
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})

		Convey("Domain does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    "/api/v1/domains/0",
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.DomainNotExist)
			So(err.Message, ShouldEqual, "Domain does not exist")
		})

		Reset(func() {
			deleteUser(u1)
			deleteUser(u2)
			deleteToken(t1)
			deleteToken(t2)
			deleteDomain(d1)
		})
	})
}

func domainRenewURL(domain *models.Domain) string {
	return domainURL(domain) + "/renew"
}

func TestAPIv1DomainRenew(t *testing.T) {
	Convey("API v1 - Domain renew", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()
		activateUser(u1)
		_, d1 := createDomain1(t1)

		Convey("Success", func() {
			d1.ExpiredAt = time.Now().AddDate(0, 0, 7)
			models.DB.Save(d1)

			time.Sleep(time.Second)

			var domain models.Domain
			r := Request(RequestOptions{
				Method: "POST",
				URL:    domainRenewURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &domain)
			So(domain.Id, ShouldEqual, d1.Id)
			So(domain.Name, ShouldEqual, d1.Name)
			So(domain.UserId, ShouldEqual, d1.UserId)
			So(domain.CreatedAt, ShouldResemble, d1.CreatedAt)
			So(domain.UpdatedAt, ShouldHappenAfter, d1.UpdatedAt)
			So(domain.ExpiredAt.Unix(), ShouldResemble, d1.ExpiredAt.AddDate(1, 0, 0).Unix())
		})

		Convey("Domain is not renewable", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    domainRenewURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusForbidden)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.DomainNotRenewable)
			So(err.Message, ShouldEqual, "This domain can not be renew until "+models.ISOTime(d1.ExpiredAt.AddDate(0, 0, -30)))
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    domainRenewURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t2.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusForbidden)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.DomainForbidden)
			So(err.Message, ShouldEqual, "You are forbidden to access this domain")
		})

		Convey("Unauthorized", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    domainRenewURL(d1),
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})

		Convey("Domain does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/domains/0/renew",
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.DomainNotExist)
			So(err.Message, ShouldEqual, "Domain does not exist")
		})

		Reset(func() {
			deleteUser(u1)
			deleteUser(u2)
			deleteToken(t1)
			deleteToken(t2)
			deleteDomain(d1)
		})
	})
}
