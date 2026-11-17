package handlers

import (
	"net/http"
	"time"

	"aman/internal/domains/app"
	"aman/internal/ports/outbound"

	"github.com/gin-gonic/gin"
)

// AppHandler is the primary (inbound) HTTP adapter for app management (CAS).
type AppHandler struct {
	appRepo outbound.IAppRepository
}

func NewAppHandler(appRepo outbound.IAppRepository) *AppHandler {
	return &AppHandler{appRepo: appRepo}
}

// Store godoc
// @Summary Store new app (CAS)
// @Tags App
// @Accept json
// @Produce json
// @Param payload body createAppBody true "payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/apps [post]
func (h AppHandler) Store(c *gin.Context) {
	var payload createAppBody
	if err := c.BindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	a, err := h.appRepo.Create(app.App{
		Name:           payload.Name,
		Url:            payload.Url,
		AuthUrl:        payload.AuthUrl,
		OtpType:        payload.OtpType,
		SessionTimeout: time.Duration(payload.SessionTimeout) * time.Second,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   a,
	})
}

// Single godoc
// @Summary Get app detail (CAS)
// @Tags App
// @Accept json
// @Produce json
// @Param id path string true "app id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/apps/{id} [get]
func (h AppHandler) Single(c *gin.Context) {
	appID := c.Param("id")

	a, err := h.appRepo.FindByID(appID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "app not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   a,
	})
}

// Index godoc
// @Summary Get apps list (CAS)
// @Tags App
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/apps [get]
func (h AppHandler) Index(c *gin.Context) {
	apps, err := h.appRepo.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   apps,
	})
}

// Update godoc
// @Summary Update app (CAS)
// @Tags App
// @Accept json
// @Produce json
// @Param id path string true "app id"
// @Param payload body updateAppBody true "payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/apps/{id} [put]
func (h AppHandler) Update(c *gin.Context) {
	appID := c.Param("id")

	var payload updateAppBody
	if err := c.BindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	existing, err := h.appRepo.FindByID(appID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "app not found.",
		})
		return
	}

	existing.Name = payload.Name
	existing.Url = payload.Url
	existing.AuthUrl = payload.AuthUrl
	existing.OtpType = payload.OtpType
	existing.SessionTimeout = time.Duration(payload.SessionTimeout) * time.Second

	updated, err := h.appRepo.Update(*existing)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   updated,
	})
}

// Delete godoc
// @Summary Delete app (CAS)
// @Tags App
// @Accept json
// @Produce json
// @Param id path string true "app id"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/apps/{id} [delete]
func (h AppHandler) Delete(c *gin.Context) {
	if err := h.appRepo.Delete(c.Param("id")); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   true,
	})
}
