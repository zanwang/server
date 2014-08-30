package models

import (
	"regexp"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/coopernurse/gorp"
	"github.com/tommy351/maji.moe/errors"
)

const (
	domainExpiry = time.Hour * 24 * 365
)

var (
	rDomainName     = regexp.MustCompile("^[a-zA-Z]+[a-zA-Z\\d\\-]*$")
	reservedDomains = []string{"www", "api", "email", "static"}
)

// Domain model
type Domain struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	ExpiredAt time.Time `db:"expired_at" json:"expired_at"`
	UserID    int64     `db:"user_id" json:"user_id"`
}

func (data *Domain) Validate(s gorp.SqlExecutor) error {
	if govalidator.IsNull(data.Name) {
		return errors.New("name", errors.Required, "Name is required")
	}

	if len(data.Name) > 63 {
		return errors.New("name", errors.MaxLength, "Maximum length of name is 63")
	}

	if !rDomainName.MatchString(data.Name) {
		return errors.New("name", errors.DomainName, "Domain name is invalid")
	}

	inarr := false

	for _, str := range reservedDomains {
		if data.Name == str {
			inarr = true
			break
		}
	}

	if inarr {
		return errors.New("name", errors.DomainReserved, "Domain name has been reserved")
	}

	if id, err := s.SelectInt("SELECT id FROM domains WHERE name=?", data.Name); err == nil {
		if data.ID != id {
			return errors.New("name", errors.DomainUsed, "Domain name has been taken")
		}
	}

	return nil
}

func (data *Domain) PreInsert(s gorp.SqlExecutor) error {
	if err := data.Validate(s); err != nil {
		return err
	}

	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	data.ExpiredAt = now.Add(domainExpiry)
	return nil
}

func (data *Domain) PreUpdate(s gorp.SqlExecutor) error {
	if err := data.Validate(s); err != nil {
		return err
	}

	data.UpdatedAt = time.Now()
	return nil
}

func (data *Domain) PreDelete(s gorp.SqlExecutor) error {
	if _, err := s.Exec("DELETE FROM records WHERE domain_id=?", data.ID); err != nil {
		return err
	}

	return nil
}

func (data *Domain) Renew() {
	data.ExpiredAt = data.ExpiredAt.Add(domainExpiry)
}
