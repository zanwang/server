package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
	"github.com/tommy351/maji.moe/util"
)

func (a *APIv1) RecordList(c *gin.Context) {
	var records []models.Record
	domain := c.MustGet("domain").(*models.Domain)

	if _, err := models.DB.Select(&records, "SELECT * FROM records WHERE domain_id=?", domain.ID); err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusOK, records)
}

type recordForm struct {
	Name     *string `json:"name"`
	Type     *string `json:"type"`
	Value    *string `json:"value"`
	TTL      *uint   `json:"ttl"`
	Priority *uint   `json:"priority"`
}

func (f *recordForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&f.Name:     "name",
		&f.Type:     "type",
		&f.Value:    "value",
		&f.TTL:      "ttl",
		&f.Priority: "priority",
	}
}

func (a *APIv1) RecordCreate(c *gin.Context) {
	var form recordForm

	if err := binding.Bind(c.Request, &form); err != nil {
		bindingError(err)
	}

	if form.Name == nil {
		panic(errors.New("name", errors.Required, "Name is required"))
	}

	if form.Type == nil {
		panic(errors.New("type", errors.Required, "Type is required"))
	}

	domain := c.MustGet("domain").(*models.Domain)
	record := models.Record{
		Name:     *form.Name,
		Type:     *form.Type,
		DomainID: domain.ID,
	}

	if form.Value == nil {
		record.Value = ""
	} else {
		record.Value = *form.Value
	}

	if form.TTL == nil {
		record.TTL = 0
	} else {
		record.TTL = *form.TTL
	}

	if form.Priority == nil {
		record.Priority = 0
	} else {
		record.Priority = *form.Priority
	}

	if err := models.DB.Insert(&record); err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusCreated, record)
}

func (a *APIv1) RecordShow(c *gin.Context) {
	record := c.MustGet("record").(*models.Record)

	util.Render.JSON(c.Writer, http.StatusOK, record)
}

func (a *APIv1) RecordUpdate(c *gin.Context) {
	var form recordForm

	if err := binding.Bind(c.Request, &form); err != nil {
		bindingError(err)
	}

	record := c.MustGet("record").(*models.Record)

	if form.Name != nil {
		record.Name = *form.Name
	}

	if form.Type != nil {
		record.Type = *form.Type
	}

	if form.Value != nil {
		record.Value = *form.Value
	}

	if form.TTL != nil {
		record.TTL = *form.TTL
	}

	if form.Priority != nil {
		record.Priority = *form.Priority
	}

	if _, err := models.DB.Update(record); err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusOK, record)
}

func (a *APIv1) RecordDestroy(c *gin.Context) {
	record := c.MustGet("record").(*models.Record)

	if _, err := models.DB.Delete(record); err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
