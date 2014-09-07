package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/models"
	"github.com/majimoe/server/util"
	"github.com/mholt/binding"
)

func (a *APIv1) RecordList(c *gin.Context) {
	var records []models.Record
	domain := c.MustGet("domain").(*models.Domain)

	if err := models.DB.Where("domain_id = ?", domain.Id).Find(&records).Error; err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusOK, records)
}

type recordForm struct {
	Name     *string `json:"name"`
	Type     *string `json:"type"`
	Value    *string `json:"value"`
	TTL      *int    `json:"ttl"`
	Priority *int    `json:"priority"`
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

	if form.Type == nil {
		panic(errors.New("type", errors.Required, "Type is required"))
	}

	domain := c.MustGet("domain").(*models.Domain)
	record := models.Record{
		Type:     *form.Type,
		DomainId: domain.Id,
	}

	if form.Name != nil {
		record.Name = *form.Name
	}

	if form.Value != nil {
		record.Value = *form.Value
	}

	if form.TTL != nil {
		record.Ttl = *form.TTL
	}

	if form.Priority != nil {
		record.Priority = *form.Priority
	}

	if err := models.DB.Create(&record).Error; err != nil {
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
		record.Ttl = *form.TTL
	}

	if form.Priority != nil {
		record.Priority = *form.Priority
	}

	if err := models.DB.Save(record).Error; err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusOK, record)
}

func (a *APIv1) RecordDestroy(c *gin.Context) {
	record := c.MustGet("record").(*models.Record)

	if err := models.DB.Delete(record).Error; err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
