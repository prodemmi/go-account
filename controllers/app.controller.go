package controllers

import (
	"account/services"
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
	Name string `form:"name" json:"name" binding:"required"`
	Url  string `form:"url" json:"url" binding:"required"`
}

func (a AppController) Store(context *gin.Context) {
	var payload CreateAppBody
	if err := context.BindJSON(&payload); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if apps, err := a.AppService.CreateApp(payload.Name, payload.Url); err == nil {
		context.JSON(http.StatusOK, gin.H{
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

func (a AppController) Single(context *gin.Context) {
	appID := context.Param("app")

	if len(appID) != 0 {
		app, err := a.AppService.GetApp(appID)

		if app == nil {
			context.JSON(http.StatusNotFound, gin.H{
				"error": "app not found.",
			})
			context.Abort()
			return
		}

		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"status": "error",
				"data":   app,
			})
			context.Abort()
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   app,
		})
	}

}

func (a AppController) Index(context *gin.Context) {
	apps, err := a.AppService.GetAppsList()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   apps,
	})
}

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
