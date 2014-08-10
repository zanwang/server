package controllers

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/tommy351/maji.moe/models"
)

func DomainList(params martini.Params, token *models.Token, r render.Render, db *gorp.DbMap, user *models.User) {
	/*var user models.User

	if err := dbMap.SelectOne(&user, "SELECT id FROM users WHERE id=?", params["user_id"]); err != nil {
		r.Status(http.StatusNotFound)
		return
	}*/

	var domains []models.Domain
	var err error

	if user.ID == token.UserID {
		_, err = db.Select(&domains, "SELECT * FROM domains WHERE user_id=?", user.ID)
	} else {
		_, err = db.Select(&domains, "SELECT * FROM domains WHERE user_id=? AND public=?", user.ID, true)
	}

	if err != nil {
		panic(err)
	}

	r.JSON(http.StatusOK, domains)
}

type DomainCreateForm struct {
	Name   string `form:"name"`
	Public bool   `form:"public"`
}

func (form *DomainCreateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Name, "name").Required("").Length(2, 63, "")

	return errors
}

func DomainCreate(form DomainCreateForm, db *gorp.DbMap, r render.Render, user *models.User, token *models.Token) {
	var domain models.Domain

	if err := db.SelectOne(&domain, "SELECT id FROM domains WHERE name=?", form.Name); err != nil {
		errors := newErr([]string{"name"}, "212", "Domain name has been taken")
		r.JSON(http.StatusBadRequest, FormatErr(errors))
		return
	}

	domain = models.Domain{
		Name:   form.Name,
		Public: form.Public,
		UserID: token.UserID,
	}

	if err := db.Insert(&domain); err != nil {
		panic(err)
	}

	r.JSON(http.StatusCreated, domain)
}

func DomainShow(r render.Render, domain *models.Domain) {
	/*
		var domain models.Domain

		if err := dbMap.SelectOne(&domain, "SELECT * FROM domains WHERE id=?", params["domain_id"]); err != nil {
			r.Status(http.StatusNotFound)
			return
		}*/
	/*
		if domain.UserID != token.UserID && !domain.Public {
			r.Status(http.StatusForbidden)
			return
		}
	*/
	r.JSON(http.StatusOK, domain)
}

type DomainUpdateForm struct {
	Name   string `form:"name"`
	Public bool   `form:"public"`
}

func (form *DomainUpdateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Name, "name").Length(2, 63, "")

	return errors
}

func DomainUpdate(form DomainUpdateForm, r render.Render, db *gorp.DbMap, domain *models.Domain) {
	/*
		var domain models.Domain

		if err := dbMap.SelectOne(&domain, "SELECT * FROM domains WHERE id=?", params["domain_id"]); err != nil {
			r.Status(http.StatusNotFound)
			return
		}
	*/
	/*
		if domain.UserID != token.UserID {
			r.Status(http.StatusForbidden)
			return
		}
	*/
	if form.Name != "" {
		domain.Name = form.Name
	}

	domain.Public = form.Public

	if count, err := db.Update(domain); count > 0 {
		r.JSON(http.StatusOK, domain)
	} else if err != nil {
		panic(err)
	} else {
		r.Status(http.StatusNotFound)
	}
}

func DomainDestroy(res http.ResponseWriter, db *gorp.DbMap, domain *models.Domain) {
	/*
		var domain models.Domain

		if err := dbMap.SelectOne(&domain, "SELECT id,user_id FROM domains WHERE id=?", params["domain_id"]); err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}*/
	/*
		if domain.UserID != token.UserID {
			res.WriteHeader(http.StatusForbidden)
			return
		}
	*/
	if count, err := db.Delete(domain); count > 0 {
		res.WriteHeader(http.StatusNoContent)
	} else if err != nil {
		panic(err)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}
