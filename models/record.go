package models

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/majimoe/server/errors"
)

var (
	rDomain    = regexp.MustCompile("\\.[a-zA-Z]{2,}$")
	RecordType = []string{"A", "CNAME", "MX", "TXT", "SPF", "AAAA", "NS", "LOC"}
)

type Record struct {
	Id        int64
	Name      string
	Type      string
	Value     string
	Ttl       int `sql:"ttl"`
	Priority  int
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DomainId  int64     `json:"domain_id"`
}

func (r Record) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":         r.Id,
		"name":       r.Name,
		"type":       r.Type,
		"value":      r.Value,
		"ttl":        r.Ttl,
		"priority":   r.Priority,
		"created_at": ISOTime(r.CreatedAt),
		"updated_at": ISOTime(r.UpdatedAt),
		"domain_id":  r.DomainId,
	})
}

func (r *Record) BeforeSave() error {
	if len(r.Name) > 63 {
		return errors.New("name", errors.MaxLength, "Maximum length of name is 63")
	}

	if len(r.Name) > 0 && !rDomainName.MatchString(r.Name) {
		return errors.New("name", errors.DomainName, "Only numbers and characters are allowed in domain name")
	}

	if govalidator.IsNull(r.Type) {
		return errors.New("type", errors.Required, "Type is required")
	}

	r.Type = strings.ToUpper(r.Type)
	inarr := false

	for _, str := range RecordType {
		if r.Type == str {
			inarr = true
			break
		}
	}

	if !inarr {
		return errors.New("type", errors.RecordType, "Type must be one of "+strings.Join(RecordType, ", "))
	}

	if govalidator.IsNull(r.Value) {
		return errors.New("value", errors.Required, "Value is required")
	}

	if r.Ttl < 0 {
		return errors.New("ttl", errors.Min, "Minimum value of TTL is 0")
	}

	if r.Priority < 0 {
		return errors.New("priority", errors.Min, "Minimum value of priority is 0")
	}

	if r.Ttl != 0 && (r.Ttl > 86400 || r.Ttl < 300) {
		return errors.New("ttl", errors.Range, "TTL must be between 300-86400 seconds")
	}

	if r.Priority > 65535 {
		return errors.New("priority", errors.Max, "Maximum value of priority is 65535")
	}

	switch r.Type {
	case "A":
		if !govalidator.IsIP(r.Value, 4) {
			return errors.New("value", errors.IPv4, "Value is not a valid IPv4")
		}
	case "AAAA":
		if !govalidator.IsIP(r.Value, 6) {
			return errors.New("value", errors.IPv6, "Value is not a valid IPv6")
		}
	case "CNAME", "MX", "NS":
		if !rDomain.MatchString(r.Value) {
			return errors.New("value", errors.Domain, "Value is not a valid domain")
		}
	}

	r.UpdatedAt = time.Now().UTC()

	return nil
}

func (r *Record) BeforeCreate() error {
	r.CreatedAt = time.Now().UTC()
	return nil
}
