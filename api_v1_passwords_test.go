package main

import (
	"net/http"
	"testing"

	"github.com/majimoe/server/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAPIv1PasswordReset(t *testing.T) {
	_, user := createUser1()

	defer func() {
		deleteUser(user)
	}()

	Convey("API v1 - Password reset", t, func() {
		Convey("Success", func() {
			var data map[string]interface{}
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/passwords/reset",
				Body: map[string]interface{}{
					"email": user.Email,
				},
			})

			So(r.Code, ShouldEqual, http.StatusAccepted)
			ParseJSON(r.Body, &data)
			So(data["email"], ShouldResemble, user.Email)
		})

		Convey("Email is required", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/passwords/reset",
				Body:   map[string]interface{}{},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "email")
			So(err.Code, ShouldEqual, errors.Required)
			So(err.Message, ShouldEqual, "Email is required")
		})

		Convey("Email is invalid", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/passwords/reset",
				Body: map[string]interface{}{
					"email": "abc",
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
				URL:    "/api/v1/passwords/reset",
				Body: map[string]interface{}{
					"email": "abc@def.com",
				},
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.UserNotExist)
			So(err.Message, ShouldEqual, "User does not exist")
		})
	})
}
