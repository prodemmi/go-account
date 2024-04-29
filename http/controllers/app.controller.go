package controllers

import (
	"account/database/models"
	"account/http/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	AppController struct {
		AppService services.IAppService
	}
	IAppController interface {
		Store(context *gin.Context)
		Single(context *gin.Context)
		Index(context *gin.Context)
		Delete(context *gin.Context)
	}
)

type CreateAppBody struct {
	Name    string `form:"name" json:"name" binding:"required" example:"fotomi"`
	Url     string `form:"url" json:"url" binding:"required" example:"https://fotomi.net/"`
	AuthUrl string `form:"auth_url" json:"auth_url" binding:"required" example:"https://fotomi.net/auth/callback"`
	OtpType string `json:"otp_type" binding:"required" example:"email"`
}

// Store
// @Summary Store new server
// @Tags App
// @Accept json
// @Produce json
// @Param payload body CreateAppBody true "payload"
// @Success 200 {object} swagger.SuccessJSONResult{data=models.App}
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/apps [post]
func (a AppController) Store(context *gin.Context) {
	var payload CreateAppBody
	if err := context.BindJSON(&payload); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if apps, err := a.AppService.CreateApp(models.App{
		Name:    payload.Name,
		Url:     payload.Url,
		AuthUrl: payload.AuthUrl,
		OtpType: payload.OtpType,
	}); err == nil {
		context.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"data":   apps,
		})
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
	}
}

// Single
// @Summary Get server detail
// @Tags App
// @Accept json
// @Produce json
// @Param id query string true "server id"
// @Success 200 {object} swagger.SuccessJSONResult{data=models.App}
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/apps/:id [get]
func (a AppController) Single(context *gin.Context) {
	appID := context.Param("id")

	if len(appID) != 0 {
		app, err := a.AppService.GetApp(appID)

		if app == nil {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "app not found.",
			})
			return
		}

		if err != nil {
			context.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status": "error",
				"data":   app,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   app,
		})
	}

}

// Index
// @Summary Get apps detail
// @Tags App
// @Accept json
// @Produce json
// @Success 200 {object} swagger.SuccessJSONResult{data=[]models.App}
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/apps/ [get]
func (a AppController) Index(context *gin.Context) {
	apps, err := a.AppService.GetAppsList()

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   apps,
	})
}

// Delete
// @Summary Delete server
// @Tags App
// @Accept json
// @Produce json
// @Success 200 {object} swagger.SuccessJSONResult{data=boolean}
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/apps/:id [delete]
func (a AppController) Delete(context *gin.Context) {
	apps, err := a.AppService.DeleteApp(context.Param("id"))

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   apps,
	})
}

func NewAppController(appService services.IAppService) AppController {
	return AppController{
		appService,
	}
}
