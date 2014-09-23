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

func createToken(body map[string]interface{}) (*httptest.ResponseRecorder, *models.Token) {
	var token models.Token
	r := Request(RequestOptions{
		Method: "POST",
		URL:    "/api/v1/tokens",
		Body:   body,
	})

	if err := ParseJSON(r.Body, &token); err != nil {
		panic(err)
	}

	return r, &token
}

func createToken1() (*httptest.ResponseRecorder, *models.Token) {
	return createToken(map[string]interface{}{
		"email":    Fixture.Users[0].Email,
		"password": Fixture.Users[0].Password,
	})
}

func createToken2() (*httptest.ResponseRecorder, *models.Token) {
	return createToken(map[string]interface{}{
		"email":    Fixture.Users[1].Email,
		"password": Fixture.Users[1].Password,
	})
}

func deleteToken(token *models.Token) {
	if err := models.DB.Delete(token).Error; err != nil {
		panic(err)
	}
}

func TestAPIv1TokenCreate(t *testing.T) {
	_, user := createUser1()

	defer func() {
		deleteUser(user)
	}()

	Convey("API v1 - Token create", t, func() {
		Convey("Success", func() {
			r, token := createToken1()
			defer deleteToken(token)

			So(r.Code, ShouldEqual, http.StatusCreated)
			So(r.Header().Get("Pragma"), ShouldEqual, "no-cache")
			So(r.Header().Get("Cache-Control"), ShouldEqual, "no-cache, no-store, must-revalidate")
			So(r.Header().Get("Expires"), ShouldEqual, "0")
			So(len(token.Key), ShouldEqual, 64)
			So(token.UpdatedAt.AddDate(0, 0, 7), ShouldResemble, token.ExpiredAt)
			So(token.UserId, ShouldEqual, user.Id)
		})

		Convey("Duplicated token", func() {
			_, t1 := createToken1()
			time.Sleep(time.Second)
			r, t2 := createToken1()
			defer deleteToken(t1)
			defer deleteToken(t2)

			So(r.Code, ShouldEqual, http.StatusCreated)
			So(t2.UserId, ShouldEqual, t1.UserId)
			So(t2.Key, ShouldNotEqual, t1.Key)
			So(t2.UpdatedAt, ShouldHappenAfter, t1.UpdatedAt)
		})

		Convey("Email required", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/tokens",
				Body:   map[string]interface{}{},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "email")
			So(err.Code, ShouldEqual, errors.Required)
			So(err.Message, ShouldEqual, "Email is required")
		})

		Convey("Password required", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/tokens",
				Body: map[string]interface{}{
					"email": "abc@def.com",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "password")
			So(err.Code, ShouldEqual, errors.Required)
			So(err.Message, ShouldEqual, "Password is required")
		})

		Convey("Invalid email", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/tokens",
				Body: map[string]interface{}{
					"email":    "abc",
					"password": "123456",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "email")
			So(err.Code, ShouldEqual, errors.Email)
			So(err.Message, ShouldEqual, "Email is invalid")
		})

		Convey("User does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/tokens",
				Body: map[string]interface{}{
					"email":    "abc@def.com",
					"password": "123456",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "email")
			So(err.Code, ShouldEqual, errors.UserNotExist)
			So(err.Message, ShouldEqual, "User does not exist")
		})

		Convey("Password too short", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/tokens",
				Body: map[string]interface{}{
					"email":    Fixture.Users[0].Email,
					"password": "123",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "password")
			So(err.Code, ShouldEqual, errors.MinLength)
			So(err.Message, ShouldEqual, "Minimum length of password is 6")
		})

		Convey("Password too long", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/tokens",
				Body: map[string]interface{}{
					"email":    Fixture.Users[0].Email,
					"password": "123465465431354354354343436546531354343543544654654",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "password")
			So(err.Code, ShouldEqual, errors.MaxLength)
			So(err.Message, ShouldEqual, "Maximum length of password is 50")
		})

		Convey("Wrong password", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/tokens",
				Body: map[string]interface{}{
					"email":    Fixture.Users[0].Email,
					"password": "rqweqflafksdpof",
				},
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "password")
			So(err.Code, ShouldEqual, errors.WrongPassword)
			So(err.Message, ShouldEqual, "Password is wrong")
		})

		Convey("Password unset", func() {
			user.Password = ""
			models.DB.Save(user)

			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/tokens",
				Body: map[string]interface{}{
					"email":    Fixture.Users[0].Email,
					"password": "rqweqflafksdpof",
				},
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "password")
			So(err.Code, ShouldEqual, errors.PasswordUnset)
			So(err.Message, ShouldEqual, "Password has not been set")
		})
	})
}

func TestAPIv1TokenList(t *testing.T) {
	_, u1 := createUser1()
	_, u2 := createUser2()
	_, t1a := createToken1()
	_, t1b := createToken1()
	_, t1c := createToken1()
	_, t2 := createToken2()

	defer func() {
		deleteUser(u1)
		deleteUser(u2)
		deleteToken(t1a)
		deleteToken(t1b)
		deleteToken(t1c)
		deleteToken(t2)
	}()

	Convey("API v1 - Token list", t, func() {
		Convey("Success", func() {
			var tokens []models.Token
			r := Request(RequestOptions{
				Method: "GET",
				URL:    "/api/v1/users/" + strconv.FormatInt(u1.Id, 10) + "/tokens",
				Headers: map[string]string{
					"Authorization": "token " + t1a.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			So(r.Header().Get("Pragma"), ShouldEqual, "no-cache")
			So(r.Header().Get("Cache-Control"), ShouldEqual, "no-cache, no-store, must-revalidate")
			So(r.Header().Get("Expires"), ShouldEqual, "0")
			ParseJSON(r.Body, &tokens)
			So(len(tokens), ShouldEqual, 3)
			So(tokens[0].IsCurrent, ShouldBeTrue)
			So(tokens[0].Ip.String(), ShouldEqual, "10.0.0.1")
			So(tokens[1].IsCurrent, ShouldBeFalse)
			So(tokens[1].Ip.String(), ShouldEqual, "10.0.0.1")
			So(tokens[2].IsCurrent, ShouldBeFalse)
			So(tokens[2].Ip.String(), ShouldEqual, "10.0.0.1")
		})

		Convey("Unauthorized", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "GET",
				URL:    "/api/v1/users/" + strconv.FormatInt(u1.Id, 10) + "/tokens",
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "GET",
				URL:    "/api/v1/users/" + strconv.FormatInt(u1.Id, 10) + "/tokens",
				Headers: map[string]string{
					"Authorization": "token " + t2.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusForbidden)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.UserForbidden)
			So(err.Message, ShouldEqual, "You are forbidden to access this user")
		})
	})
}

func TestAPIv1TokenUpdate(t *testing.T) {
	_, user := createUser1()
	_, token := createToken1()

	defer func() {
		deleteUser(user)
		deleteToken(token)
	}()

	Convey("API v1 - Token update", t, func() {
		Convey("Success", func() {
			// Wait for a while to let token updated
			time.Sleep(time.Second)

			var t models.Token
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    "/api/v1/tokens",
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			So(r.Header().Get("Pragma"), ShouldEqual, "no-cache")
			So(r.Header().Get("Cache-Control"), ShouldEqual, "no-cache, no-store, must-revalidate")
			So(r.Header().Get("Expires"), ShouldEqual, "0")
			ParseJSON(r.Body, &t)
			So(t.Key, ShouldEqual, token.Key)
			So(t.UpdatedAt, ShouldHappenAfter, token.UpdatedAt)
			So(t.UpdatedAt.AddDate(0, 0, 7).Unix(), ShouldEqual, t.ExpiredAt.Unix())
			So(t.UserId, ShouldEqual, user.Id)
		})

		Convey("Unauthorized", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    "/api/v1/tokens",
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})
	})
}

func TestAPIv1TokenDestroy(t *testing.T) {
	_, user := createUser1()

	defer func() {
		deleteUser(user)
	}()

	Convey("API v1 - Token destroy", t, func() {
		_, token := createToken1()

		Reset(func() {
			deleteToken(token)
		})

		Convey("Success", func() {
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    "/api/v1/tokens",
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNoContent)

			// Confirm token has been deleted
			var count int

			if err := models.DB.Table("tokens").Where("`key` = ?", token.Key).Count(&count).Error; err != nil {
				panic(err)
			}

			So(count, ShouldEqual, 0)
		})

		Convey("Unauthorized", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    "/api/v1/tokens",
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})
	})
}
