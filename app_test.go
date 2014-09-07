package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/majimoe/server/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserActivation(t *testing.T) {
	Convey("User activation", t, func() {
		_, user := createUser1()

		Convey("Success", func() {
			// Get complete user data from database
			models.DB.First(user, user.Id)

			r := Request(RequestOptions{
				Method: "GET",
				URL:    fmt.Sprintf("/users/%d/activation/%s", user.Id, user.ActivationToken),
			})

			So(r.Code, ShouldEqual, http.StatusFound)
			So(r.Header().Get("Location"), ShouldEqual, "/app")
		})

		Convey("User does not exist", func() {
			r := Request(RequestOptions{
				Method: "GET",
				URL:    fmt.Sprintf("/users/%d/activation/%s", 0, "abcdefghijklmnopqrstuvwxyzabcdef"),
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("Wrong token", func() {
			r := Request(RequestOptions{
				Method: "GET",
				URL:    fmt.Sprintf("/users/%d/activation/%s", user.Id, "abcdefghijklmnopqrstuvwxyzabcdef"),
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("User has been activated", func() {
			// Get complete user data from database
			models.DB.First(user, user.Id)
			activateUser(user)

			r := Request(RequestOptions{
				Method: "GET",
				URL:    fmt.Sprintf("/users/%d/activation/%s", user.Id, user.ActivationToken),
			})

			So(r.Code, ShouldEqual, http.StatusFound)
			So(r.Header().Get("Location"), ShouldEqual, "/app")
		})

		Reset(func() {
			deleteUser(user)
		})
	})
}
