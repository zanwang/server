package controllers

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/models"
	"github.com/majimoe/server/util"
	"github.com/mholt/binding"
)

type tokenForm struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (f *tokenForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&f.Email:    "email",
		&f.Password: "password",
	}
}

func noCacheHeader(c *gin.Context) {
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Expires", "0")
}

func responseToken(c *gin.Context, status int, token *models.Token) {
	noCacheHeader(c)
	util.Render.JSON(c.Writer, status, token)
}

func (a *APIv1) TokenList(c *gin.Context) {
	tokens := []models.Token{}
	user := c.MustGet("user").(*models.User)
	token := c.MustGet("token").(*models.Token)

	if err := models.DB.Where("user_id = ?", user.Id).Find(&tokens).Error; err != nil {
		if err != gorm.RecordNotFound {
			panic(err)
		}
	}

	result := make([]map[string]interface{}, len(tokens))

	for i, t := range tokens {
		result[i] = map[string]interface{}{
			"key":        t.Key,
			"updated_at": models.ISOTime(t.UpdatedAt),
			"expired_at": models.ISOTime(t.GetExpiredTime()),
			"user_id":    t.UserId,
			"ip":         t.Ip.String(),
			"is_current": t.Key == token.Key,
		}
	}

	noCacheHeader(c)
	util.Render.JSON(c.Writer, http.StatusOK, result)
}

func (a *APIv1) TokenCreate(c *gin.Context) {
	var form tokenForm

	if err := binding.Bind(c.Request, &form); err != nil {
		bindingError(err)
	}

	if form.Email == nil {
		panic(errors.New("email", errors.Required, "Email is required"))
	}

	if form.Password == nil {
		panic(errors.New("password", errors.Required, "Password is required"))
	}

	if !govalidator.IsEmail(*form.Email) {
		panic(errors.New("email", errors.Email, "Email is invalid"))
	}

	var user models.User

	if err := models.DB.Where("email = ?", *form.Email).Find(&user).Error; err != nil {
		panic(&errors.API{
			Status:  http.StatusBadRequest,
			Field:   "email",
			Code:    errors.UserNotExist,
			Message: "User does not exist",
		})
	}

	if err := user.Authenticate(*form.Password); err != nil {
		panic(err)
	}

	token := models.Token{
		UserId: user.Id,
	}

	token.SetIP(GetIPFromContext(c))

	if err := models.DB.Create(&token).Error; err != nil {
		panic(err)
	}

	responseToken(c, http.StatusCreated, &token)
}

func (a *APIv1) TokenUpdate(c *gin.Context) {
	token := c.MustGet("token").(*models.Token)
	token.SetIP(GetIPFromContext(c))

	if err := models.DB.Save(token).Error; err != nil {
		panic(err)
	}

	responseToken(c, http.StatusOK, token)
}

func (a *APIv1) TokenDestroy(c *gin.Context) {
	token := c.MustGet("token").(*models.Token)
	defer c.Set("token", nil)

	if err := models.DB.Delete(token).Error; err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
