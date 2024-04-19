package controllers

import (
	"account/database/models"
	"account/http/services"
	"account/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type (
	UserController struct {
		UserService services.IUserService
	}

	IUserController interface {
		Index(context *gin.Context)
		Store(context *gin.Context)
		Single(context *gin.Context)
	}
)

// Index
// @Summary Get users detail
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} swagger.SuccessJSONResult{data=[]models.User}
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/users/ [get]
func (c UserController) Index(context *gin.Context) {
	appID := context.GetHeader("APP-ID")
	users, err := c.UserService.GetUsersList(appID)

	if err == nil {
		context.AbortWithStatusJSON(http.StatusOK, gin.H{
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
	ReagentID    *string `json:"reagent_id" example:"6513270e-269e-4d37-b2a7-4de452e6b438"`
	FirstName    string  `json:"first_name" example:"Emad"`
	LastName     string  `json:"last_name" example:"Malekpour"`
	Email        string  `json:"email" example:"prodemmi@gmail.com"`
	Mobile       string  `json:"mobile" example:"+989130243519"`
	NationalCode *string `json:"national_code" example:"1279856451"`
}

// Store
// @Summary Store new user
// @Tags User
// @Accept json
// @Produce json
// @Param APP-ID header string true "null"
// @Param payload body StoreUserBody true "payload"
// @Success 200 {object} swagger.SuccessJSONResult{data=models.User}
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/users [post]
func (c UserController) Store(context *gin.Context) {
	var payload StoreUserBody
	appID := context.GetHeader("APP-ID")

	if err := context.BindJSON(&payload); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	app, _ := uuid.Parse(appID)
	createUser := models.User{
		AppID:        app,
		Email:        &payload.Email,
		Mobile:       &payload.Mobile,
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		NationalCode: payload.NationalCode,
		SignedUp:     false,
	}

	user, err := c.UserService.Create(createUser)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   user,
	})
}

// Single
// @Summary Get user detail
// @Tags User
// @Accept json
// @Produce json
// @Param id query string true "user id"
// @Success 200 {object} swagger.SuccessJSONResult{data=models.User}
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/users/:id [get]
func (c UserController) Single(context *gin.Context) {
	id := context.Param("id")

	if !utils.IsUUID(id) {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "id is not a valid uuid.",
		})
		return
	}

	user, err := c.UserService.GetUser(id)

	if user == nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "user not found.",
		})
		return
	}

	if err == nil {
		context.AbortWithStatusJSON(http.StatusOK, gin.H{
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

func NewUserController(userService services.IUserService) IUserController {
	return UserController{
		UserService: userService,
	}
}
