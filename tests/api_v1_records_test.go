package tests

import (
	"net/http"
	"strconv"
	"strings"

	. "github.com/onsi/gomega"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
)

func (s *TestSuite) APIv1Record() {
	s.Describe("Record", func() {
		s.APIv1RecordCreate()
		s.APIv1RecordList()
		s.APIv1RecordShow()
	})
}

func recordCreateURL(id int64) string {
	return "/api/v1/domains/" + strconv.FormatInt(id, 10) + "/records"
}

func recordURL(id int64) string {
	return "/api/v1/records/" + strconv.FormatInt(id, 10)
}

func (s *TestSuite) createRecord(key string, token *models.Token, domain *models.Domain, body map[string]interface{}) {
	var record models.Record
	r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
		Body: body,
		Headers: map[string]string{
			"Authorization": "token " + token.Key,
		},
	})

	Expect(r.Code).To(Equal(http.StatusCreated))

	s.ParseJSON(r.Body, &record)
	Expect(record.Name).To(Equal(body["name"]))
	Expect(record.Type).To(Equal(body["type"]))
	Expect(record.Value).To(Equal(body["value"]))
	Expect(record.TTL).To(Equal(body["ttl"]))
	Expect(record.Priority).To(Equal(body["priority"]))
	Expect(record.DomainID).To(Equal(domain.ID))

	s.Set(key, &record)
}

func (s *TestSuite) deleteRecord(key string) {
	record := s.Get(key).(*models.Record)
	models.DB.Delete(record)
	s.Del(key)
}

func (s *TestSuite) createRecord1() {
	token := s.Get("token").(*models.Token)
	domain := s.Get("domain").(*models.Domain)
	s.createRecord("record", token, domain, map[string]interface{}{
		"name":     Fixture.Records[0].Name,
		"type":     Fixture.Records[0].Type,
		"value":    Fixture.Records[0].Value,
		"ttl":      Fixture.Records[0].TTL,
		"priority": Fixture.Records[0].Priority,
	})
}

func (s *TestSuite) deleteRecord1() {
	s.deleteRecord("record")
}

func (s *TestSuite) createRecord2() {
	token := s.Get("token2").(*models.Token)
	domain := s.Get("domain2").(*models.Domain)
	s.createRecord("record2", token, domain, map[string]interface{}{
		"name":     Fixture.Records[1].Name,
		"type":     Fixture.Records[1].Type,
		"value":    Fixture.Records[1].Value,
		"ttl":      Fixture.Records[1].TTL,
		"priority": Fixture.Records[1].Priority,
	})
}

func (s *TestSuite) deleteRecord2() {
	s.deleteRecord("record2")
}

func (s *TestSuite) APIv1RecordCreate() {
	s.Describe("Create", func() {
		s.Before(func() {
			s.createUser1()
			s.createUser2()
			s.createToken1()
			s.createToken2()
			s.createDomain1()
		})

		s.It("Success", func() {
			s.createRecord1()
		})

		s.It("Name started with number", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name":  "1a2b",
					"type":  "A",
					"value": "127.0.0.1",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("name"))
			Expect(err.Code).To(Equal(errors.DomainName))
			Expect(err.Message).To(Equal("Domain name is invalid"))
		})

		s.It("Name with special characters", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name":  "中文",
					"type":  "A",
					"value": "127.0.0.1",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("name"))
			Expect(err.Code).To(Equal(errors.DomainName))
			Expect(err.Message).To(Equal("Domain name is invalid"))
		})

		s.It("Name too long", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name":  "erfwoerjweijrliejrwejrliwejrwjerliwwjeroiljweloirjweolirireorweorjweorwoerjwoerwoeirlsj",
					"type":  "A",
					"value": "127.0.0.1",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("name"))
			Expect(err.Code).To(Equal(errors.MaxLength))
			Expect(err.Message).To(Equal("Maximum length of name is 63"))
		})

		s.It("Type is required", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"name":  "1a2b",
					"value": "127.0.0.1",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("type"))
			Expect(err.Code).To(Equal(errors.Required))
			Expect(err.Message).To(Equal("Type is required"))
		})

		s.It("Wrong type", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"type":  "WTF",
					"value": "127.0.0.1",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("type"))
			Expect(err.Code).To(Equal(errors.RecordType))
			Expect(err.Message).To(Equal("Type must be one of " + strings.Join(models.RecordType, ", ")))
		})

		s.It("Transform type to uppercase", func() {
			var record models.Record
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "a",
					"value": "127.0.0.1",
				},
			})

			Expect(r.Code).To(Equal(http.StatusCreated))

			s.ParseJSON(r.Body, &record)
			Expect(record.Type).To(Equal("A"))

			models.DB.Delete(&record)
		})

		s.It("Value is required", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]string{
					"type": "A",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("value"))
			Expect(err.Code).To(Equal(errors.Required))
			Expect(err.Message).To(Equal("Value is required"))
		})

		s.It("TTL too small", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "127.0.0.1",
					"ttl":   10,
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("ttl"))
			Expect(err.Code).To(Equal(errors.Range))
			Expect(err.Message).To(Equal("TTL must be between 300-86400 seconds"))
		})

		s.It("TTL too large", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "127.0.0.1",
					"ttl":   100000,
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("ttl"))
			Expect(err.Code).To(Equal(errors.Range))
			Expect(err.Message).To(Equal("TTL must be between 300-86400 seconds"))
		})

		s.It("Negative TTL", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "127.0.0.1",
					"ttl":   -10,
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("ttl"))
			Expect(err.Code).To(Equal(errors.Min))
			Expect(err.Message).To(Equal("Minimum value of TTL is 0"))
		})

		s.It("Negative priority", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":     "MX",
					"value":    "mail.google.com",
					"priority": -10,
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("priority"))
			Expect(err.Code).To(Equal(errors.Min))
			Expect(err.Message).To(Equal("Minimum value of priority is 0"))
		})

		s.It("Priority too large", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":     "MX",
					"value":    "mail.google.com",
					"priority": 100000,
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("priority"))
			Expect(err.Code).To(Equal(errors.Max))
			Expect(err.Message).To(Equal("Maximum value of priority is 65535"))
		})

		s.It("A test - success", func() {
			var record models.Record
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "127.0.0.1",
				},
			})

			Expect(r.Code).To(Equal(http.StatusCreated))

			s.ParseJSON(r.Body, &record)
			Expect(record.Name).To(Equal(""))
			Expect(record.Type).To(Equal("A"))
			Expect(record.Value).To(Equal("127.0.0.1"))
			Expect(record.TTL).To(BeEquivalentTo(0))
			Expect(record.Priority).To(BeEquivalentTo(0))
			Expect(record.DomainID).To(Equal(domain.ID))

			models.DB.Delete(&record)
		})

		s.It("A test - fail", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "wtwerowjeroiw",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("value"))
			Expect(err.Code).To(Equal(errors.IPv4))
			Expect(err.Message).To(Equal("Value is not a valid IPv4"))
		})

		s.It("AAAA test - success", func() {
			var record models.Record
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "AAAA",
					"value": "2404:6800:4008:c04::8a",
				},
			})

			Expect(r.Code).To(Equal(http.StatusCreated))

			s.ParseJSON(r.Body, &record)
			Expect(record.Name).To(Equal(""))
			Expect(record.Type).To(Equal("AAAA"))
			Expect(record.Value).To(Equal("2404:6800:4008:c04::8a"))
			Expect(record.TTL).To(BeEquivalentTo(0))
			Expect(record.Priority).To(BeEquivalentTo(0))
			Expect(record.DomainID).To(Equal(domain.ID))

			models.DB.Delete(&record)
		})

		s.It("AAAA test - fail", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "AAAA",
					"value": "wtwerowjeroiw",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("value"))
			Expect(err.Code).To(Equal(errors.IPv6))
			Expect(err.Message).To(Equal("Value is not a valid IPv6"))
		})

		s.It("CNAME test - success", func() {
			var record models.Record
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "CNAME",
					"value": "github.com",
				},
			})

			Expect(r.Code).To(Equal(http.StatusCreated))

			s.ParseJSON(r.Body, &record)
			Expect(record.Name).To(Equal(""))
			Expect(record.Type).To(Equal("CNAME"))
			Expect(record.Value).To(Equal("github.com"))
			Expect(record.TTL).To(BeEquivalentTo(0))
			Expect(record.Priority).To(BeEquivalentTo(0))
			Expect(record.DomainID).To(Equal(domain.ID))

			models.DB.Delete(&record)
		})

		s.It("CNAME test - fail", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "CNAME",
					"value": "wtwerowjeroiw",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("value"))
			Expect(err.Code).To(Equal(errors.Domain))
			Expect(err.Message).To(Equal("Value is not a valid domain"))
		})

		s.It("MX test - success", func() {
			var record models.Record
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":     "MX",
					"value":    "mail.google.com",
					"priority": 10,
				},
			})

			Expect(r.Code).To(Equal(http.StatusCreated))

			s.ParseJSON(r.Body, &record)
			Expect(record.Name).To(Equal(""))
			Expect(record.Type).To(Equal("MX"))
			Expect(record.Value).To(Equal("mail.google.com"))
			Expect(record.TTL).To(BeEquivalentTo(0))
			Expect(record.Priority).To(BeEquivalentTo(10))
			Expect(record.DomainID).To(Equal(domain.ID))

			models.DB.Delete(&record)
		})

		s.It("MX test - fail", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "MX",
					"value": "wtwerowjeroiw",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("value"))
			Expect(err.Code).To(Equal(errors.Domain))
			Expect(err.Message).To(Equal("Value is not a valid domain"))
		})

		s.It("NS test - success", func() {
			var record models.Record
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "NS",
					"value": "pam.ns.cloudflare.com",
				},
			})

			Expect(r.Code).To(Equal(http.StatusCreated))

			s.ParseJSON(r.Body, &record)
			Expect(record.Name).To(Equal(""))
			Expect(record.Type).To(Equal("NS"))
			Expect(record.Value).To(Equal("pam.ns.cloudflare.com"))
			Expect(record.TTL).To(BeEquivalentTo(0))
			Expect(record.Priority).To(BeEquivalentTo(0))
			Expect(record.DomainID).To(Equal(domain.ID))

			models.DB.Delete(&record)
		})

		s.It("NS test - fail", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "NS",
					"value": "wtwerowjeroiw",
				},
			})

			Expect(r.Code).To(Equal(http.StatusBadRequest))

			s.ParseJSON(r.Body, &err)
			Expect(err.Field).To(Equal("value"))
			Expect(err.Code).To(Equal(errors.Domain))
			Expect(err.Message).To(Equal("Value is not a valid domain"))
		})

		s.It("Forbidden (with wrong token)", func() {
			var err errors.API
			token := s.Get("token2").(*models.Token)
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "127.0.0.1",
				},
			})

			Expect(r.Code).To(Equal(http.StatusForbidden))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.DomainForbidden))
			Expect(err.Message).To(Equal("You are forbidden to access this domain"))
		})

		s.It("Unauthorized (without token)", func() {
			var err errors.API
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("POST", recordCreateURL(domain.ID), &requestOptions{
				Body: map[string]interface{}{
					"type":  "A",
					"value": "127.0.0.1",
				},
			})

			Expect(r.Code).To(Equal(http.StatusUnauthorized))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.TokenRequired))
			Expect(err.Message).To(Equal("Token is required"))
		})

		s.It("Domain does not exist", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			r := s.Request("POST", recordCreateURL(999999999), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
				Body: map[string]interface{}{
					"type":  "A",
					"value": "127.0.0.1",
				},
			})

			Expect(r.Code).To(Equal(http.StatusNotFound))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.DomainNotExist))
			Expect(err.Message).To(Equal("Domain does not exist"))
		})

		s.After(func() {
			s.deleteUser1()
			s.deleteUser2()
			s.deleteToken1()
			s.deleteToken2()
			s.deleteDomain1()
			s.deleteRecord1()
		})
	})
}

func (s *TestSuite) APIv1RecordList() {
	s.Describe("List", func() {
		s.Before(func() {
			s.createUser1()
			s.createUser2()
			s.createToken1()
			s.createToken2()
			s.createDomain1()
			s.createDomain2()
			s.createRecord1()
			s.createRecord2()
		})

		s.It("Success", func() {
			var records []models.Record
			domain := s.Get("domain").(*models.Domain)
			record := s.Get("record").(*models.Record)
			token := s.Get("token").(*models.Token)
			r := s.Request("GET", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusOK))

			s.ParseJSON(r.Body, &records)
			Expect(records).To(HaveLen(1))
			re := records[0]
			Expect(re.ID).To(Equal(record.ID))
			Expect(re.Name).To(Equal(record.Name))
			Expect(re.Type).To(Equal(record.Type))
			Expect(re.Value).To(Equal(record.Value))
			Expect(re.TTL).To(Equal(record.TTL))
			Expect(re.Priority).To(Equal(record.Priority))
			Expect(re.DomainID).To(Equal(record.DomainID))
		})

		s.It("Forbidden (with wrong token)", func() {
			var err errors.API
			domain := s.Get("domain").(*models.Domain)
			token := s.Get("token2").(*models.Token)
			r := s.Request("GET", recordCreateURL(domain.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusForbidden))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.DomainForbidden))
			Expect(err.Message).To(Equal("You are forbidden to access this domain"))
		})

		s.It("Unauthorized (without token)", func() {
			var err errors.API
			domain := s.Get("domain").(*models.Domain)
			r := s.Request("GET", recordCreateURL(domain.ID), nil)

			Expect(r.Code).To(Equal(http.StatusUnauthorized))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.TokenRequired))
			Expect(err.Message).To(Equal("Token is required"))
		})

		s.It("Domain does not exist", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			r := s.Request("GET", recordCreateURL(999999999), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusNotFound))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.DomainNotExist))
			Expect(err.Message).To(Equal("Domain does not exist"))
		})

		s.After(func() {
			s.deleteUser1()
			s.deleteUser2()
			s.deleteToken1()
			s.deleteToken2()
			s.deleteDomain1()
			s.deleteDomain2()
			s.deleteRecord1()
			s.deleteRecord2()
		})
	})
}

func (s *TestSuite) APIv1RecordShow() {
	s.Describe("Show", func() {
		s.Before(func() {
			s.createUser1()
			s.createUser2()
			s.createToken1()
			s.createToken2()
			s.createDomain1()
			s.createRecord1()
		})

		s.It("Success", func() {
			var re map[string]interface{}
			record := s.Get("record").(*models.Record)
			token := s.Get("token").(*models.Token)
			r := s.Request("GET", recordURL(record.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusOK))

			s.ParseJSON(r.Body, &re)
			Expect(re["id"]).To(BeEquivalentTo(record.ID))
			Expect(re["name"]).To(Equal(record.Name))
			Expect(re["type"]).To(Equal(record.Type))
			Expect(re["value"]).To(Equal(record.Value))
			Expect(re["ttl"]).To(BeEquivalentTo(record.TTL))
			Expect(re["priority"]).To(BeEquivalentTo(record.Priority))
			Expect(re["domain_id"]).To(BeEquivalentTo(record.DomainID))
			Expect(re).To(HaveKey("created_at"))
			Expect(re).To(HaveKey("updated_at"))
		})

		s.It("Forbidden (with wrong token)", func() {
			var err errors.API
			record := s.Get("record").(*models.Record)
			token := s.Get("token2").(*models.Token)
			r := s.Request("GET", recordURL(record.ID), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusForbidden))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.RecordForbidden))
			Expect(err.Message).To(Equal("You are forbidden to access this record"))
		})

		s.It("Unauthorized (without token)", func() {
			var err errors.API
			record := s.Get("record").(*models.Record)
			r := s.Request("GET", recordURL(record.ID), nil)

			Expect(r.Code).To(Equal(http.StatusUnauthorized))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.TokenRequired))
			Expect(err.Message).To(Equal("Token is required"))
		})

		s.It("Record does not exist", func() {
			var err errors.API
			token := s.Get("token").(*models.Token)
			r := s.Request("GET", recordURL(9999999999), &requestOptions{
				Headers: map[string]string{
					"Authorization": "token " + token.Key,
				},
			})

			Expect(r.Code).To(Equal(http.StatusNotFound))

			s.ParseJSON(r.Body, &err)
			Expect(err.Code).To(Equal(errors.RecordNotExist))
			Expect(err.Message).To(Equal("Record does not exist"))
		})

		s.After(func() {
			s.deleteUser1()
			s.deleteUser2()
			s.deleteToken1()
			s.deleteToken2()
			s.deleteDomain1()
			s.deleteRecord1()
		})
	})
}
