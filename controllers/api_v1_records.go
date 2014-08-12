package controllers

import (
	"net/http"
	"strings"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/tommy351/maji.moe/models"
)

func RecordList(r render.Render, db *gorp.DbMap, domain *models.Domain) {
	var records []models.Record

	if _, err := db.Select(&records, "SELECT * FROM records WHERE domain_id=?", domain.ID); err != nil {
		panic(err)
	}

	r.JSON(http.StatusOK, records)
}

type RecordCreateForm struct {
	Name     string `form:"name"`
	Type     string `form:"type"`
	Content  string `form:"content"`
	TTL      uint   `form:"ttl"`
	Priority uint   `form:"priority"`
}

func (form *RecordCreateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	form.Type = strings.ToUpper(form.Type)

	if form.Name != "" {
		v.Validate(&form.Name, "name").Length(1, 63, "")
	}

	v.Validate(&form.Type, "type").Required("").IsIn(models.RecordType, "")
	v.Validate(&form.Content, "content").Required("")
	v.Validate(&form.TTL, "ttl").Required("").Within(120, 86400, "")

	if form.Type == "MX" {
		v.Validate(&form.Priority, "priority").Required("").Min(0, "")
	}

	switch form.Type {
	case "A":
		v.Validate(&form.Content, "content").IP("")
	case "CNAME", "MX", "NS":
		v.Validate(&form.Content, "content").Domain("")
	case "AAAA":
		v.Validate(&form.Content, "content").IPv6("")
	}

	return errors
}

func RecordCreate(form RecordCreateForm, db *gorp.DbMap, r render.Render, domain *models.Domain) {
	record := models.Record{
		Name:     form.Name,
		Type:     form.Type,
		Content:  form.Content,
		TTL:      form.TTL,
		Priority: form.Priority,
		DomainID: domain.ID,
	}

	if err := db.Insert(&record); err != nil {
		panic(err)
	}

	r.JSON(http.StatusCreated, record)
}

func RecordShow(r render.Render, record *models.Record) {
	r.JSON(http.StatusOK, record)
}

type RecordUpdateForm struct {
	Name     string `form:"name"`
	Type     string `form:"type"`
	Content  string `form:"content"`
	TTL      uint   `form:"ttl"`
	Priority uint   `form:"priority"`
}

func (form *RecordUpdateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	form.Type = strings.ToUpper(form.Type)

	if form.Name != "" {
		v.Validate(&form.Name, "name").Length(1, 63, "")
	}

	v.Validate(&form.Type, "type").IsIn(models.RecordType, "")
	v.Validate(&form.TTL, "ttl").Within(120, 86400, "")

	if form.Type == "MX" {
		v.Validate(&form.Priority, "priority").Min(0, "")
	}

	switch form.Type {
	case "A":
		v.Validate(&form.Content, "content").IP("")
	case "CNAME", "MX", "NS":
		v.Validate(&form.Content, "content").Domain("")
	case "AAAA":
		v.Validate(&form.Content, "content").IPv6("")
	}

	return errors
}

func RecordUpdate(form RecordUpdateForm, db *gorp.DbMap, r render.Render, record *models.Record) {
	if form.Name != "" {
		record.Name = form.Name
	}

	if form.Type != "" {
		record.Type = form.Type
	}

	if form.Content != "" {
		record.Content = form.Content
	}

	record.TTL = form.TTL
	record.Priority = form.Priority

	if count, err := db.Update(record); count > 0 {
		r.JSON(http.StatusOK, record)
	} else {
		panic(err)
	}
}

func RecordDestroy(db *gorp.DbMap, res http.ResponseWriter, record *models.Record) {
	if count, err := db.Delete(record); count > 0 {
		res.WriteHeader(http.StatusNoContent)
	} else if err != nil {
		panic(err)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}
