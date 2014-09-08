package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/models"
	. "github.com/smartystreets/goconvey/convey"
)

func createUser(body map[string]interface{}) (*httptest.ResponseRecorder, *models.User) {
	var user models.User
	r := Request(RequestOptions{
		Method: "POST",
		URL:    "/api/v1/users",
		Body:   body,
	})

	if err := ParseJSON(r.Body, &user); err != nil {
		panic(err)
	}

	return r, &user
}

func createUser1() (*httptest.ResponseRecorder, *models.User) {
	return createUser(map[string]interface{}{
		"name":     Fixture.Users[0].Name,
		"email":    Fixture.Users[0].Email,
		"password": Fixture.Users[0].Password,
	})
}

func createUser2() (*httptest.ResponseRecorder, *models.User) {
	return createUser(map[string]interface{}{
		"name":     Fixture.Users[1].Name,
		"email":    Fixture.Users[1].Email,
		"password": Fixture.Users[1].Password,
	})
}

func md5str(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

func deleteUser(user *models.User) {
	if err := models.DB.Delete(user).Error; err != nil {
		panic(err)
	}
}

func activateUser(user *models.User) {
	user.Activated = true

	if err := models.DB.Save(user).Error; err != nil {
		panic(err)
	}
}

func TestAPIv1UserCreate(t *testing.T) {
	Convey("API v1 - User create", t, func() {
		Convey("Success", func() {
			r, user := createUser1()
			defer deleteUser(user)

			So(r.Code, ShouldEqual, http.StatusCreated)
			So(user.Name, ShouldEqual, Fixture.Users[0].Name)
			So(user.Email, ShouldEqual, Fixture.Users[0].Email)
			So(user.Activated, ShouldBeFalse)
			So(user.Password, ShouldBeBlank)
			So(user.Avatar, ShouldEqual, "//www.gravatar.com/avatar/"+md5str(Fixture.Users[0].Email))
			So(user.CreatedAt, ShouldNotBeNil)
			So(user.UpdatedAt, ShouldNotBeNil)
			So(user.Password, ShouldBeBlank)
			So(user.ActivationToken, ShouldBeBlank)
			So(user.PasswordResetToken, ShouldBeBlank)
			So(user.FacebookId, ShouldBeBlank)
			So(user.GoogleId, ShouldBeBlank)
		})

		Convey("Name required", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/users",
				Body:   map[string]interface{}{},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "name")
			So(err.Code, ShouldEqual, errors.Required)
			So(err.Message, ShouldEqual, "Name is required")
		})

		Convey("Password required", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/users",
				Body: map[string]interface{}{
					"name": "John",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "password")
			So(err.Code, ShouldEqual, errors.Required)
			So(err.Message, ShouldEqual, "Password is required")
		})

		Convey("Email required", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/users",
				Body: map[string]interface{}{
					"name":     "John",
					"password": "123456",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "email")
			So(err.Code, ShouldEqual, errors.Required)
			So(err.Message, ShouldEqual, "Email is required")
		})

		Convey("Password too short", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/users",
				Body: map[string]interface{}{
					"name":     "John",
					"password": "123",
					"email":    "john@maji.moe",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "password")
			So(err.Code, ShouldEqual, errors.Length)
			So(err.Message, ShouldEqual, "The length of password must be between 6-50")
		})

		Convey("Password too long", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/users",
				Body: map[string]interface{}{
					"name":     "John",
					"password": "123456465465465465465645456456456456456456456456456456456",
					"email":    "john@maji.moe",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "password")
			So(err.Code, ShouldEqual, errors.Length)
			So(err.Message, ShouldEqual, "The length of password must be between 6-50")
		})

		Convey("Email format", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/users",
				Body: map[string]interface{}{
					"name":     "John",
					"password": "123456",
					"email":    "john",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "email")
			So(err.Code, ShouldEqual, errors.Email)
			So(err.Message, ShouldEqual, "Email is invalid")
		})

		Convey("Email has been taken", func() {
			_, user := createUser1()
			defer deleteUser(user)

			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    "/api/v1/users",
				Body: map[string]interface{}{
					"name":     Fixture.Users[0].Name,
					"password": Fixture.Users[0].Password,
					"email":    Fixture.Users[0].Email,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "email")
			So(err.Code, ShouldEqual, errors.EmailUsed)
			So(err.Message, ShouldEqual, "Email has been taken")
		})
	})
}

func userShowURL(user *models.User) string {
	return "/api/v1/users/" + strconv.FormatInt(user.Id, 10)
}

func TestAPIv1UserShow(t *testing.T) {
	Convey("API v1 - User show", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()

		Reset(func() {
			deleteUser(u1)
			deleteUser(u2)
			deleteToken(t1)
			deleteToken(t2)
		})

		Convey("Success (private)", func() {
			var user models.User
			r := Request(RequestOptions{
				Method: "GET",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &user)
			So(user.Id, ShouldEqual, u1.Id)
			So(user.Name, ShouldEqual, u1.Name)
			So(user.Email, ShouldEqual, u1.Email)
			So(user.Avatar, ShouldEqual, u1.Avatar)
			So(user.Activated, ShouldEqual, u1.Activated)
			So(user.CreatedAt, ShouldResemble, u1.CreatedAt)
			So(user.UpdatedAt, ShouldResemble, u1.UpdatedAt)
			So(user.Password, ShouldBeBlank)
			So(user.ActivationToken, ShouldBeBlank)
			So(user.PasswordResetToken, ShouldBeBlank)
			So(user.FacebookId, ShouldBeBlank)
			So(user.GoogleId, ShouldBeBlank)
		})

		Convey("Success (public)", func() {
			var user models.User
			r := Request(RequestOptions{
				Method: "GET",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t2.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &user)
			So(user.Id, ShouldEqual, u1.Id)
			So(user.Name, ShouldEqual, u1.Name)
			So(user.Email, ShouldBeBlank)
			So(user.Avatar, ShouldEqual, u1.Avatar)
			So(user.Activated, ShouldBeFalse)
			So(user.CreatedAt, ShouldResemble, u1.CreatedAt)
			So(user.UpdatedAt, ShouldResemble, u1.UpdatedAt)
			So(user.Password, ShouldBeBlank)
			So(user.ActivationToken, ShouldBeBlank)
			So(user.PasswordResetToken, ShouldBeBlank)
			So(user.FacebookId, ShouldBeBlank)
			So(user.GoogleId, ShouldBeBlank)
		})

		Convey("Success (Unauthorized)", func() {
			var user models.User
			r := Request(RequestOptions{
				Method: "GET",
				URL:    userShowURL(u1),
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &user)
			So(user.Id, ShouldEqual, u1.Id)
			So(user.Name, ShouldEqual, u1.Name)
			So(user.Email, ShouldBeBlank)
			So(user.Avatar, ShouldEqual, u1.Avatar)
			So(user.Activated, ShouldBeFalse)
			So(user.CreatedAt, ShouldResemble, u1.CreatedAt)
			So(user.UpdatedAt, ShouldResemble, u1.UpdatedAt)
			So(user.Password, ShouldBeBlank)
			So(user.ActivationToken, ShouldBeBlank)
			So(user.PasswordResetToken, ShouldBeBlank)
			So(user.FacebookId, ShouldBeBlank)
			So(user.GoogleId, ShouldBeBlank)
		})

		Convey("User does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "GET",
				URL:    "/api/v1/users/0",
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.UserNotExist)
			So(err.Message, ShouldEqual, "User does not exist")
		})
	})
}

func TestAPIv1UserUpdate(t *testing.T) {
	Convey("API v1 - User update", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()

		Reset(func() {
			deleteUser(u1)
			deleteUser(u2)
			deleteToken(t1)
			deleteToken(t2)
		})

		Convey("Success", func() {
			time.Sleep(time.Second)

			var user models.User
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]string{
					"name": "WTF",
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &user)
			So(user.Id, ShouldEqual, u1.Id)
			So(user.Name, ShouldEqual, "WTF")
			So(user.Email, ShouldEqual, u1.Email)
			So(user.Avatar, ShouldEqual, u1.Avatar)
			So(user.Activated, ShouldEqual, u1.Activated)
			So(user.CreatedAt, ShouldResemble, u1.CreatedAt)
			So(user.UpdatedAt, ShouldHappenAfter, u1.UpdatedAt)
			So(user.Password, ShouldBeBlank)
			So(user.ActivationToken, ShouldBeBlank)
			So(user.PasswordResetToken, ShouldBeBlank)
			So(user.FacebookId, ShouldBeBlank)
			So(user.GoogleId, ShouldBeBlank)
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
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
				Method: "PUT",
				URL:    userShowURL(u1),
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})

		Convey("Edit email", func() {
			activateUser(u1)

			var user models.User
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"email": "abc@def.com",
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &user)
			So(user.Email, ShouldEqual, "abc@def.com")
			So(user.Avatar, ShouldEqual, "//www.gravatar.com/avatar/"+md5str(user.Email))
			So(user.Activated, ShouldBeFalse)
		})

		Convey("Invalid email", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
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

		Convey("Email has been taken", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"email": Fixture.Users[1].Email,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "email")
			So(err.Code, ShouldEqual, errors.EmailUsed)
			So(err.Message, ShouldEqual, "Email has been taken")
		})

		Convey("Email wasn't changed", func() {
			activateUser(u1)

			var user models.User
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"email": Fixture.Users[0].Email,
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &user)
			So(user.Email, ShouldEqual, u1.Email)
			So(user.Activated, ShouldBeTrue)
		})

		Convey("Password too short", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"old_password": Fixture.Users[0].Password,
					"password":     "123",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "password")
			So(err.Code, ShouldEqual, errors.Length)
			So(err.Message, ShouldEqual, "The length of password must be between 6-50")
		})

		Convey("Password too long", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"old_password": Fixture.Users[0].Password,
					"password":     "123464646546546546546546546546546546546546546546546546546546546",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "password")
			So(err.Code, ShouldEqual, errors.Length)
			So(err.Message, ShouldEqual, "The length of password must be between 6-50")
		})

		Convey("Wrong current password", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"old_password": "wejroiwjerlfsadsd",
					"password":     "1234567",
				},
			})

			So(r.Code, ShouldEqual, http.StatusForbidden)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "old_password")
			So(err.Code, ShouldEqual, errors.WrongPassword)
			So(err.Message, ShouldEqual, "Password is wrong")
		})

		Convey("Current password is required", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"password": "1234567",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "old_password")
			So(err.Code, ShouldEqual, errors.Required)
			So(err.Message, ShouldEqual, "Current password is required")
		})

		Convey("Modify password", func() {
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"old_password": Fixture.Users[0].Password,
					"password":     "12345678",
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)

			// Confirm password has been changed
			var user models.User

			if err := models.DB.First(&user, u1.Id).Error; err != nil {
				panic(err)
			}

			So(user.Authenticate("12345678"), ShouldBeNil)
		})

		Convey("Don't need current password if password isn't set before", func() {
			u1.Password = ""
			models.DB.Save(u1)

			r := Request(RequestOptions{
				Method: "PUT",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"password": "12345678",
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)

			// Confirm password has been changed
			var user models.User

			if err := models.DB.First(&user, u1.Id).Error; err != nil {
				panic(err)
			}

			So(user.Authenticate("12345678"), ShouldBeNil)
		})

		Convey("User does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    "/api/v1/users/0",
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.UserNotExist)
			So(err.Message, ShouldEqual, "User does not exist")
		})
	})
}

func TestAPIv1UserDestroy(t *testing.T) {
	Convey("API v1 - User destroy", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()

		Reset(func() {
			deleteUser(u1)
			deleteUser(u2)
			deleteToken(t1)
			deleteToken(t2)
		})

		Convey("Success", func() {
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    userShowURL(u1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNoContent)

			// Confirm user has been deleted
			var count int

			if err := models.DB.Table("users").Where("id = ?", u1.Id).Count(&count).Error; err != nil {
				panic(err)
			}

			So(count, ShouldEqual, 0)
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    userShowURL(u1),
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
				Method: "DELETE",
				URL:    userShowURL(u1),
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})

		Convey("User does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    "/api/v1/users/0",
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.UserNotExist)
			So(err.Message, ShouldEqual, "User does not exist")
		})
	})
}
