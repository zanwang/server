package models

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/majimoe/server/errors"
)

const (
	domainExpiry      = time.Hour * 24 * 365
	domainRenewPeriod = time.Hour * 24 * 30
)

var (
	rDomainName     = regexp.MustCompile("^[a-zA-Z\\d\\-]+$")
	reservedDomains = []string{"www", "api", "email", "static", "test"}
)

type Domain struct {
	Id        int64
	Name      string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiredAt time.Time `json:"expired_at"`
	UserId    int64     `json:"user_id"`
}

func (d Domain) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":         d.Id,
		"name":       d.Name,
		"created_at": ISOTime(d.CreatedAt),
		"updated_at": ISOTime(d.UpdatedAt),
		"expired_at": ISOTime(d.ExpiredAt),
		"user_id":    d.UserId,
	})
}

func (d *Domain) BeforeSave() error {
	if d.Name == "" {
		return errors.New("name", errors.Required, "Name is required")
	}

	if len(d.Name) > 63 {
		return errors.New("name", errors.MaxLength, "Maximum length of name is 63")
	}

	if !rDomainName.MatchString(d.Name) {
		return errors.New("name", errors.DomainName, "Only numbers and characters are allowed in domain name")
	}

	inarr := false

	for _, str := range reservedDomains {
		if d.Name == str {
			inarr = true
			break
		}
	}

	if inarr {
		return errors.New("name", errors.DomainReserved, "Domain name has been reserved")
	}

	d.UpdatedAt = time.Now().UTC()

	return nil
}

func (d *Domain) BeforeCreate() error {
	now := time.Now().UTC()
	d.CreatedAt = now
	d.ExpiredAt = now.Add(domainExpiry)
	return nil
}

func (d *Domain) Renew() error {
	if time.Now().Add(domainRenewPeriod).Before(d.ExpiredAt) {
		return &errors.API{
			Status:  http.StatusForbidden,
			Code:    errors.DomainNotRenewable,
			Message: "This domain can not be renew until " + ISOTime(d.ExpiredAt.Add(-domainRenewPeriod)),
		}
	}

	d.ExpiredAt = d.ExpiredAt.Add(domainExpiry)

	return nil
}
