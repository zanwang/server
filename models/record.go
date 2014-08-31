package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/coopernurse/gorp"
	"github.com/majimoe/server/errors"
)

var rDomain = regexp.MustCompile("\\.[a-zA-Z]{2,}$")
var RecordType = []string{"A", "CNAME", "MX", "TXT", "SPF", "AAAA", "NS", "LOC"}

// Record model
type Record struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"`
	Value     string    `db:"value" json:"value"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	DomainID  int64     `db:"domain_id" json:"domain_id"`
	TTL       int       `db:"ttl" json:"ttl"`
	Priority  int       `db:"priority" json:"priority"`
}

func (data *Record) Validate() error {
	if len(data.Name) > 63 {
		return errors.New("name", errors.MaxLength, "Maximum length of name is 63")
	}

	if len(data.Name) > 0 && !rDomainName.MatchString(data.Name) {
		return errors.New("name", errors.DomainName, "Domain name is invalid")
	}

	if govalidator.IsNull(data.Type) {
		return errors.New("type", errors.Required, "Type is required")
	}

	data.Type = strings.ToUpper(data.Type)
	inarr := false

	for _, str := range RecordType {
		if data.Type == str {
			inarr = true
			break
		}
	}

	if !inarr {
		return errors.New("type", errors.RecordType, "Type must be one of "+strings.Join(RecordType, ", "))
	}

	if govalidator.IsNull(data.Value) {
		return errors.New("value", errors.Required, "Value is required")
	}

	if data.TTL < 0 {
		return errors.New("ttl", errors.Min, "Minimum value of TTL is 0")
	}

	if data.Priority < 0 {
		return errors.New("priority", errors.Min, "Minimum value of priority is 0")
	}

	if data.TTL != 0 && (data.TTL > 86400 || data.TTL < 300) {
		return errors.New("ttl", errors.Range, "TTL must be between 300-86400 seconds")
	}

	if data.Priority > 65535 {
		return errors.New("priority", errors.Max, "Maximum value of priority is 65535")
	}

	switch data.Type {
	case "A":
		if !govalidator.IsIP(data.Value, 4) {
			return errors.New("value", errors.IPv4, "Value is not a valid IPv4")
		}
	case "AAAA":
		if !govalidator.IsIP(data.Value, 6) {
			return errors.New("value", errors.IPv6, "Value is not a valid IPv6")
		}
	case "CNAME", "MX", "NS":
		if !rDomain.MatchString(data.Value) {
			return errors.New("value", errors.Domain, "Value is not a valid domain")
		}
	}

	return nil
}

func (data *Record) PreInsert(s gorp.SqlExecutor) error {
	if err := data.Validate(); err != nil {
		return err
	}

	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	return nil
}

func (data *Record) PreUpdate(s gorp.SqlExecutor) error {
	if err := data.Validate(); err != nil {
		return err
	}

	data.UpdatedAt = time.Now()
	return nil
}
