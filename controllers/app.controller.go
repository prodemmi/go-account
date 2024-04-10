package controllers

import (
	"github.com/gin-gonic/gin"
	"account/services"
	"net/http"
)

type AppController struct {
	App        *gin.Engine
	AppService services.AppService
}

type appController interface {
	Store(context *gin.Context)
	Single(context *gin.Context)
	Index(context *gin.Context)
	DeleteOneApp(context *gin.Context)
}

type CreateAppBody struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (a AppController) Store(context *gin.Context) {
	var payload CreateAppBody
	if err := context.BindJSON(&payload); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	apps, err := a.AppService.Create(payload.Name, payload.Url)
	if err == nil {
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
		app, err := a.AppService.FindOne(appID)

		if app == nil {
			context.JSON(http.StatusNotFound, gin.H{
				"error": "app not found.",
			})
			return
		}

		if err != nil {
			context.JSON(http.StatusOK, gin.H{
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

func (a AppController) Index(context *gin.Context) {
	apps, err := a.AppService.FindAll()

	if err != nil {
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   apps,
	})
}

func (a AppController) DeleteOneApp(context *gin.Context) {

}
