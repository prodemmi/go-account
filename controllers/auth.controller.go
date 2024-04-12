package controllers

import (
	"account/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	AuthService services.IAuthService
}

type IAuthController interface {
	SignUp(context *gin.Context)
	SignIn(context *gin.Context)
	Me(context *gin.Context)
	OAuth2(context *gin.Context)
	OAuth2Fallback(context *gin.Context)
}

func (c AuthController) SignUp(context *gin.Context) {}

func (c AuthController) SignIn(context *gin.Context) {}

// Me
// @Summary Get authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "null"
// @Param API-TOKEN header string true "null"
// @Success 200 {object} swagger.SuccessJSONResult{data=models.User}
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/auth/me [get]
func (c AuthController) Me(context *gin.Context) {
	user, exist := context.Get("User")
	if exist {
		context.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   user,
		})
		return
	} else {
		context.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "user not authorized.",
		})
	}

}

func (c AuthController) OAuth2(context *gin.Context) {}

func (c AuthController) OAuth2Fallback(context *gin.Context) {}

func NewAuthController(authService services.IAuthService) IAuthController {
	return AuthController{
		AuthService: authService,
	}
}
