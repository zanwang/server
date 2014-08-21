package controllers

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/tommy351/maji.moe/config"
	"github.com/tommy351/maji.moe/models"
)

func DomainList(params martini.Params, token *models.Token, r render.Render, db *gorp.DbMap, user *models.User) {
	var domains []models.Domain

	if _, err := db.Select(&domains, "SELECT * FROM domains WHERE user_id=?", user.ID); err != nil {
		panic(err)
	}

	r.JSON(http.StatusOK, domains)
}

type DomainCreateForm struct {
	Name string `form:"name" json:"name"`
}

func (form *DomainCreateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Name, "name").Required("").MaxLength(63, "").DomainName("")

	return errors
}

func DomainCreate(form DomainCreateForm, db *gorp.DbMap, r render.Render, user *models.User, token *models.Token, conf *config.Config) {
	for _, x := range conf.ReservedDomains {
		if form.Name == x {
			errors := NewErr([]string{"name"}, "212", "Domain name has been taken")
			r.JSON(http.StatusBadRequest, FormatErr(errors))
			return
		}
	}

	var domain models.Domain

	if err := db.SelectOne(&domain, "SELECT id FROM domains WHERE name=?", form.Name); err == nil {
		errors := NewErr([]string{"name"}, "212", "Domain name has been taken")
		r.JSON(http.StatusBadRequest, FormatErr(errors))
		return
	}

	domain = models.Domain{
		Name:   form.Name,
		UserID: token.UserID,
	}

	if err := db.Insert(&domain); err != nil {
		panic(err)
	}

	r.JSON(http.StatusCreated, domain)
}

func DomainShow(r render.Render, domain *models.Domain) {
	r.JSON(http.StatusOK, domain)
}

type DomainUpdateForm struct {
	Name string `form:"name" json:"name"`
}

func (form *DomainUpdateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Name, "name").MaxLength(63, "").DomainName("")

	return errors
}

func DomainUpdate(form DomainUpdateForm, r render.Render, db *gorp.DbMap, domain *models.Domain) {
	if form.Name != "" {
		domain.Name = form.Name
	}

	if count, err := db.Update(domain); count > 0 {
		r.JSON(http.StatusOK, domain)
	} else if err != nil {
		panic(err)
	} else {
		r.Status(http.StatusNotFound)
	}
}

func DomainDestroy(res http.ResponseWriter, db *gorp.DbMap, domain *models.Domain) {
	if count, err := db.Delete(domain); count > 0 {
		res.WriteHeader(http.StatusNoContent)
	} else if err != nil {
		panic(err)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}
