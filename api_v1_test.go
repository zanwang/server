package main

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAPIv1Entry(t *testing.T) {
	Convey("API v1 - Entry", t, func() {
		var data map[string]interface{}
		r := Request(RequestOptions{
			Method: "GET",
			URL:    "/api/v1",
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		ParseJSON(r.Body, &data)
		So(data, ShouldResemble, map[string]interface{}{
			"tokens":  "/api/v1/tokens",
			"users":   "/api/v1/users",
			"domains": "/api/v1/domains",
			"records": "/api/v1/records",
			"emails":  "/api/v1/emails",
		})
	})
}
