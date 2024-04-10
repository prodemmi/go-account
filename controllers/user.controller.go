package controllers

import (
	"account/models"
	"account/services"
	"account/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
	App         *gin.Engine
	DB          *gorm.DB
	UserService services.UserService
	AppService  services.AppService
}

type userController interface {
	Index(context *gin.Context)
	Store(context *gin.Context)
	Single(context *gin.Context)
}

func (c UserController) Index(context *gin.Context) {
	users, err := c.UserService.FindAll(context.GetHeader("APP-ID"))

	if err == nil {
		context.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   users,
		})
		return
	}

	context.JSON(http.StatusInternalServerError, gin.H{
		"status": "error",
		"error":  err.Error(),
	})
}

type StoreUserBody struct {
	ReagentID    *string `json:"reagent_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	Mobile       string  `json:"mobile"`
	NationalCode *string `json:"national_code"`
}

func (c UserController) Store(context *gin.Context) {
	var payload StoreUserBody
	appID := context.GetHeader("APP-ID")

	if err := context.BindJSON(&payload); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	app, _ := uuid.FromString(appID)

	createUser := models.User{
		AppID:        app,
		Email:        payload.Email,
		Mobile:       payload.Mobile,
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		NationalCode: payload.NationalCode,
		//ReagentID:    &reagentUUID,
	}

	user, err := c.UserService.Create(&createUser)
	if err == nil {
		context.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   user,
		})
		return
	}

	context.JSON(http.StatusInternalServerError, gin.H{
		"status": "error",
		"error":  err.Error(),
	})
	return
}

func (c UserController) Single(context *gin.Context) {
	id := context.Param("id")

	if !utils.IsUUID(id) {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "id is not a valid uuid.",
		})

		return
	}

	user, err := c.UserService.FindOne(id)

	if user == nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "user not found.",
		})
		return
	}

	if err == nil {
		context.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   user,
		})
		return
	}

	context.JSON(http.StatusInternalServerError, gin.H{
		"status": "error",
		"error":  err.Error(),
	})
}
