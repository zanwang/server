package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/models"
	. "github.com/smartystreets/goconvey/convey"
)

func recordCreateURL(domain *models.Domain) string {
	return "/api/v1/domains/" + strconv.FormatInt(domain.Id, 10) + "/records"
}

func createRecord(token *models.Token, domain *models.Domain, body map[string]interface{}) (*httptest.ResponseRecorder, *models.Record) {
	var record models.Record
	r := Request(RequestOptions{
		Method: "POST",
		URL:    recordCreateURL(domain),
		Headers: map[string]string{
			"Authorization": "token " + token.Key,
		},
		Body: body,
	})

	if err := ParseJSON(r.Body, &record); err != nil {
		panic(err)
	}

	return r, &record
}

func createRecord1(token *models.Token, domain *models.Domain) (*httptest.ResponseRecorder, *models.Record) {
	return createRecord(token, domain, map[string]interface{}{
		"name":     Fixture.Records[0].Name,
		"type":     Fixture.Records[0].Type,
		"value":    Fixture.Records[0].Value,
		"ttl":      Fixture.Records[0].TTL,
		"priority": Fixture.Records[0].Priority,
	})
}

func createRecord2(token *models.Token, domain *models.Domain) (*httptest.ResponseRecorder, *models.Record) {
	return createRecord(token, domain, map[string]interface{}{
		"name":     Fixture.Records[1].Name,
		"type":     Fixture.Records[1].Type,
		"value":    Fixture.Records[1].Value,
		"ttl":      Fixture.Records[1].TTL,
		"priority": Fixture.Records[1].Priority,
	})
}

func deleteRecord(record *models.Record) {
	if err := models.DB.Delete(record).Error; err != nil {
		panic(err)
	}
}

func TestAPIv1RecordCreate(t *testing.T) {
	Convey("API v1 - Record create", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()
		activateUser(u1)
		_, d1 := createDomain1(t1)

		Convey("Success", func() {
			r, record := createRecord1(t1, d1)
			defer deleteRecord(record)

			So(r.Code, ShouldEqual, http.StatusCreated)
			So(record.Name, ShouldEqual, Fixture.Records[0].Name)
			So(record.Type, ShouldEqual, Fixture.Records[0].Type)
			So(record.Value, ShouldEqual, Fixture.Records[0].Value)
			So(record.Ttl, ShouldEqual, Fixture.Records[0].TTL)
			So(record.Priority, ShouldEqual, Fixture.Records[0].Priority)
			So(record.DomainId, ShouldEqual, d1.Id)
			So(record.CreatedAt, ShouldNotBeNil)
			So(record.UpdatedAt, ShouldNotBeNil)
		})

		Convey("Name with special characters", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"name":  "中文",
					"type":  "A",
					"value": "127.0.0.1",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "name")
			So(err.Code, ShouldEqual, errors.DomainName)
			So(err.Message, ShouldEqual, "Only numbers and characters are allowed in domain name")
		})

		Convey("Name too long", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"name":  "erfwoerjweijrliejrwejrliwejrwjerliwwjeroiljweloirjweolirireorweorjweorwoerjwoerwoeirlsj",
					"type":  "A",
					"value": "127.0.0.1",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "name")
			So(err.Code, ShouldEqual, errors.MaxLength)
			So(err.Message, ShouldEqual, "Maximum length of name is 63")
		})

		Convey("Type is required", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"value": "127.0.0.1",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "type")
			So(err.Code, ShouldEqual, errors.Required)
			So(err.Message, ShouldEqual, "Type is required")
		})

		Convey("Wrong type", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "WTF",
					"value": "127.0.0.1",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "type")
			So(err.Code, ShouldEqual, errors.RecordType)
			So(err.Message, ShouldEqual, "Type must be one of "+strings.Join(models.RecordType, ", "))
		})

		Convey("Transform type to uppercase", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "a",
					"value": "127.0.0.1",
				},
			})

			So(r.Code, ShouldEqual, http.StatusCreated)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "A")

			deleteRecord(&record)
		})

		Convey("Value is required", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type": "A",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "value")
			So(err.Code, ShouldEqual, errors.Required)
			So(err.Message, ShouldEqual, "Value is required")
		})

		Convey("TTL too small", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "127.0.0.1",
					"ttl":   10,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "ttl")
			So(err.Code, ShouldEqual, errors.Range)
			So(err.Message, ShouldEqual, "TTL must be between 300-86400 seconds")
		})

		Convey("TTL too large", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "127.0.0.1",
					"ttl":   100000,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "ttl")
			So(err.Code, ShouldEqual, errors.Range)
			So(err.Message, ShouldEqual, "TTL must be between 300-86400 seconds")
		})

		Convey("Negative TTL", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "127.0.0.1",
					"ttl":   -10,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "ttl")
			So(err.Code, ShouldEqual, errors.Min)
			So(err.Message, ShouldEqual, "Minimum value of TTL is 0")
		})

		Convey("Priority too large", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":     "A",
					"value":    "127.0.0.1",
					"priority": 100000,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "priority")
			So(err.Code, ShouldEqual, errors.Max)
			So(err.Message, ShouldEqual, "Maximum value of priority is 65535")
		})

		Convey("Negative priority", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":     "A",
					"value":    "127.0.0.1",
					"priority": -10,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "priority")
			So(err.Code, ShouldEqual, errors.Min)
			So(err.Message, ShouldEqual, "Minimum value of priority is 0")
		})

		Convey("A test - success", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "8.8.8.8",
				},
			})

			So(r.Code, ShouldEqual, http.StatusCreated)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "A")
			So(record.Value, ShouldEqual, "8.8.8.8")

			deleteRecord(&record)
		})

		Convey("A test - fail", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "wtwerowjeroiw",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "value")
			So(err.Code, ShouldEqual, errors.IPv4)
			So(err.Message, ShouldEqual, "Value is not a valid IPv4")
		})

		Convey("AAAA test - success", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "AAAA",
					"value": "2404:6800:4008:c04::8a",
				},
			})

			So(r.Code, ShouldEqual, http.StatusCreated)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "AAAA")
			So(record.Value, ShouldEqual, "2404:6800:4008:c04::8a")

			deleteRecord(&record)
		})

		Convey("AAAA test - fail", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "AAAA",
					"value": "wtwerowjeroiw",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "value")
			So(err.Code, ShouldEqual, errors.IPv6)
			So(err.Message, ShouldEqual, "Value is not a valid IPv6")
		})

		Convey("CNAME test - success", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "CNAME",
					"value": "google.com",
				},
			})

			So(r.Code, ShouldEqual, http.StatusCreated)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "CNAME")
			So(record.Value, ShouldEqual, "google.com")

			deleteRecord(&record)
		})

		Convey("CNAME test - fail", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "CNAME",
					"value": "wtwerowjeroiw",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "value")
			So(err.Code, ShouldEqual, errors.Domain)
			So(err.Message, ShouldEqual, "Value is not a valid domain")
		})

		Convey("MX test - success", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":     "MX",
					"value":    "mail.google.com",
					"priority": 10,
				},
			})

			So(r.Code, ShouldEqual, http.StatusCreated)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "MX")
			So(record.Value, ShouldEqual, "mail.google.com")
			So(record.Priority, ShouldEqual, 10)

			deleteRecord(&record)
		})

		Convey("MX test - fail", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "MX",
					"value": "wtwerowjeroiw",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "value")
			So(err.Code, ShouldEqual, errors.Domain)
			So(err.Message, ShouldEqual, "Value is not a valid domain")
		})

		Convey("NS test - success", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "NS",
					"value": "pam.ns.cloudflare.com",
				},
			})

			So(r.Code, ShouldEqual, http.StatusCreated)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "NS")
			So(record.Value, ShouldEqual, "pam.ns.cloudflare.com")

			deleteRecord(&record)
		})

		Convey("NS test - fail", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "NS",
					"value": "wtwerowjeroiw",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "value")
			So(err.Code, ShouldEqual, errors.Domain)
			So(err.Message, ShouldEqual, "Value is not a valid domain")
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
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
				URL:    recordCreateURL(d1),
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
				URL:    "/api/v1/domains/0/records",
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

func TestAPIv1RecordList(t *testing.T) {
	Convey("API v1 - Record list", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()
		activateUser(u1)
		_, d1 := createDomain1(t1)
		_, r1 := createRecord1(t1, d1)

		Convey("Success", func() {
			var records []models.Record
			r := Request(RequestOptions{
				Method: "GET",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &records)
			So(len(records), ShouldEqual, 1)
			So(records[0], ShouldResemble, *r1)
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "GET",
				URL:    recordCreateURL(d1),
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
				Method: "GET",
				URL:    recordCreateURL(d1),
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})

		Convey("Domain does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "GET",
				URL:    "/api/v1/domains/0/records",
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
			deleteRecord(r1)
		})
	})
}

func recordURL(record *models.Record) string {
	return "/api/v1/records/" + strconv.FormatInt(record.Id, 10)
}

func TestAPIv1RecordShow(t *testing.T) {
	Convey("API v1 - Record show", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()
		activateUser(u1)
		_, d1 := createDomain1(t1)
		_, r1 := createRecord1(t1, d1)

		Convey("Success", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "GET",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &record)
			So(record, ShouldResemble, *r1)
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "GET",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t2.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusForbidden)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.RecordForbidden)
			So(err.Message, ShouldEqual, "You are forbidden to access this record")
		})

		Convey("Unauthorized", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "GET",
				URL:    recordURL(r1),
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})

		Convey("Record does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "GET",
				URL:    "/api/v1/records/0",
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.RecordNotExist)
			So(err.Message, ShouldEqual, "Record does not exist")
		})

		Reset(func() {
			deleteUser(u1)
			deleteUser(u2)
			deleteToken(t1)
			deleteToken(t2)
			deleteDomain(d1)
			deleteRecord(r1)
		})
	})
}

func TestAPIv1RecordUpdate(t *testing.T) {
	Convey("API v1 - Record update", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()
		activateUser(u1)
		_, d1 := createDomain1(t1)
		_, r1 := createRecord1(t1, d1)

		Convey("Success", func() {
			time.Sleep(time.Second)

			var record models.Record
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]string{
					"name": "abc",
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &record)
			So(record.Id, ShouldEqual, r1.Id)
			So(record.Name, ShouldEqual, "abc")
			So(record.Type, ShouldEqual, r1.Type)
			So(record.Value, ShouldEqual, r1.Value)
			So(record.Ttl, ShouldEqual, r1.Ttl)
			So(record.Priority, ShouldEqual, r1.Priority)
			So(record.DomainId, ShouldEqual, r1.DomainId)
			So(record.CreatedAt, ShouldResemble, r1.CreatedAt)
			So(record.UpdatedAt, ShouldHappenAfter, r1.UpdatedAt)
		})

		Convey("Name with special characters", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
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

		Convey("Name too long", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
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

		Convey("Wrong type", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type": "WTF",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "type")
			So(err.Code, ShouldEqual, errors.RecordType)
			So(err.Message, ShouldEqual, "Type must be one of "+strings.Join(models.RecordType, ", "))
		})

		Convey("Transform type to uppercase", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type": "a",
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "A")
		})

		Convey("TTL too small", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"ttl": 10,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "ttl")
			So(err.Code, ShouldEqual, errors.Range)
			So(err.Message, ShouldEqual, "TTL must be between 300-86400 seconds")
		})

		Convey("TTL too large", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"ttl": 100000,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "ttl")
			So(err.Code, ShouldEqual, errors.Range)
			So(err.Message, ShouldEqual, "TTL must be between 300-86400 seconds")
		})

		Convey("Negative TTL", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"ttl": -10,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "ttl")
			So(err.Code, ShouldEqual, errors.Min)
			So(err.Message, ShouldEqual, "Minimum value of TTL is 0")
		})

		Convey("Priority too large", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"priority": 100000,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "priority")
			So(err.Code, ShouldEqual, errors.Max)
			So(err.Message, ShouldEqual, "Maximum value of priority is 65535")
		})

		Convey("Negative priority", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"priority": -10,
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "priority")
			So(err.Code, ShouldEqual, errors.Min)
			So(err.Message, ShouldEqual, "Minimum value of priority is 0")
		})

		Convey("A test - success", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "8.8.8.8",
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "A")
			So(record.Value, ShouldEqual, "8.8.8.8")
		})

		Convey("A test - fail", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "wtwerowjeroiw",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "value")
			So(err.Code, ShouldEqual, errors.IPv4)
			So(err.Message, ShouldEqual, "Value is not a valid IPv4")
		})

		Convey("AAAA test - success", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "AAAA",
					"value": "2404:6800:4008:c04::8a",
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "AAAA")
			So(record.Value, ShouldEqual, "2404:6800:4008:c04::8a")
		})

		Convey("AAAA test - fail", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "AAAA",
					"value": "wtwerowjeroiw",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "value")
			So(err.Code, ShouldEqual, errors.IPv6)
			So(err.Message, ShouldEqual, "Value is not a valid IPv6")
		})

		Convey("CNAME test - success", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "CNAME",
					"value": "google.com",
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "CNAME")
			So(record.Value, ShouldEqual, "google.com")
		})

		Convey("CNAME test - fail", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "POST",
				URL:    recordCreateURL(d1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "CNAME",
					"value": "wtwerowjeroiw",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "value")
			So(err.Code, ShouldEqual, errors.Domain)
			So(err.Message, ShouldEqual, "Value is not a valid domain")
		})

		Convey("MX test - success", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":     "MX",
					"value":    "mail.google.com",
					"priority": 10,
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "MX")
			So(record.Value, ShouldEqual, "mail.google.com")
			So(record.Priority, ShouldEqual, 10)
		})

		Convey("MX test - fail", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "MX",
					"value": "wtwerowjeroiw",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "value")
			So(err.Code, ShouldEqual, errors.Domain)
			So(err.Message, ShouldEqual, "Value is not a valid domain")
		})

		Convey("NS test - success", func() {
			var record models.Record
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "NS",
					"value": "pam.ns.cloudflare.com",
				},
			})

			So(r.Code, ShouldEqual, http.StatusOK)
			ParseJSON(r.Body, &record)
			So(record.Type, ShouldEqual, "NS")
			So(record.Value, ShouldEqual, "pam.ns.cloudflare.com")
		})

		Convey("NS test - fail", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
				Body: map[string]interface{}{
					"type":  "NS",
					"value": "wtwerowjeroiw",
				},
			})

			So(r.Code, ShouldEqual, http.StatusBadRequest)
			ParseJSON(r.Body, &err)
			So(err.Field, ShouldEqual, "value")
			So(err.Code, ShouldEqual, errors.Domain)
			So(err.Message, ShouldEqual, "Value is not a valid domain")
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t2.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusForbidden)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.RecordForbidden)
			So(err.Message, ShouldEqual, "You are forbidden to access this record")
		})

		Convey("Unauthorized", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    recordURL(r1),
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})

		Convey("Record does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "PUT",
				URL:    "/api/v1/records/0",
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.RecordNotExist)
			So(err.Message, ShouldEqual, "Record does not exist")
		})

		Reset(func() {
			deleteUser(u1)
			deleteUser(u2)
			deleteToken(t1)
			deleteToken(t2)
			deleteDomain(d1)
			deleteRecord(r1)
		})
	})
}

func TestAPIv1RecordDestroy(t *testing.T) {
	Convey("API v1 - Record destroy", t, func() {
		_, u1 := createUser1()
		_, u2 := createUser2()
		_, t1 := createToken1()
		_, t2 := createToken2()
		activateUser(u1)
		_, d1 := createDomain1(t1)
		_, r1 := createRecord1(t1, d1)

		Convey("Success", func() {
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNoContent)

			// Confirm record has been deleted
			var count int

			if err := models.DB.Table("records").Where("id = ?", d1.Id).Count(&count).Error; err != nil {
				panic(err)
			}

			So(count, ShouldEqual, 0)
		})

		Convey("Forbidden", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    recordURL(r1),
				Headers: map[string]string{
					"Authorization": "token " + t2.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusForbidden)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.RecordForbidden)
			So(err.Message, ShouldEqual, "You are forbidden to access this record")
		})

		Convey("Unauthorized", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    recordURL(r1),
			})

			So(r.Code, ShouldEqual, http.StatusUnauthorized)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.TokenRequired)
			So(err.Message, ShouldEqual, "Token is required")
		})

		Convey("Record does not exist", func() {
			var err errors.API
			r := Request(RequestOptions{
				Method: "DELETE",
				URL:    "/api/v1/records/0",
				Headers: map[string]string{
					"Authorization": "token " + t1.Key,
				},
			})

			So(r.Code, ShouldEqual, http.StatusNotFound)
			ParseJSON(r.Body, &err)
			So(err.Code, ShouldEqual, errors.RecordNotExist)
			So(err.Message, ShouldEqual, "Record does not exist")
		})

		Reset(func() {
			deleteUser(u1)
			deleteUser(u2)
			deleteToken(t1)
			deleteToken(t2)
			deleteDomain(d1)
			deleteRecord(r1)
		})
	})
}
