package controllers

import (
	"account/database/models"
	"account/http/services"
	"account/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	AuthService services.IAuthService
	UserService services.IUserService
	OtpService  services.IOTPService
}

type IAuthController interface {
	SignUp(context *gin.Context)
	SignIn(context *gin.Context)
	Logout(context *gin.Context)
	Me(context *gin.Context)
	OAuth2(context *gin.Context)
	OAuth2Fallback(context *gin.Context)
	VerifyOtp(context *gin.Context)
}

type SignUpBody struct {
	FirstName string `json:"first_name" binding:"required" example:"Emad"`
	LastName  string `json:"last_name" binding:"required" example:"Malekpour"`
	Email     string `json:"email" example:"prodemmi@gmail.com"`
	Mobile    string `json:"mobile" example:"+989130243519"`
	Password  string `json:"password" binding:"required" example:"secure"`
}

type SignUpRes struct {
	UserId  string `json:"user_id" example:"4462ebfc-5f91-4ef0-9cfb-ac6e7687a66e"`
	OtpHash string `json:"otp_hash" example:"7a8c54bf6913a65e19c6cf24d5a2bd5c843b031f"`
}

// SignUp
// @Summary SignUp a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "null"
// @Param payload body SignUpBody true "payload"
// @Success 200 {object} swagger.SuccessJSONResult{data=SignUpRes}
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/auth/sign-up [post]
func (c AuthController) SignUp(context *gin.Context) {
	var payload SignUpBody
	appID, _ := context.Get("AppID")

	if err := context.ShouldBindJSON(&payload); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	appUUID, _ := uuid.Parse(fmt.Sprintf("%v", appID))
	passwordHash, _ := utils.HashPassword(payload.Password)
	user := models.User{
		AppID:     appUUID,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     &payload.Email,
		Mobile:    &payload.Mobile,
		SignedUp:  false,
		Password:  passwordHash,
	}

	user, errSignUp := c.AuthService.SignUp(user)
	if errSignUp != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errSignUp.Error(),
		})
		return
	}

	otpType, _ := context.Get("AppOtpType")
	otpHash := utils.CreateAuthHash(user, context.Request.UserAgent())
	_, errOtp := c.OtpService.CreateOTP(user.ID, user.AppID, fmt.Sprintf("%v", otpType), otpHash)
	if errOtp != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errOtp.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user_id":  user.ID,
			"otp_hash": otpHash,
		},
	})
}

type SignInBody struct {
	Email    string `json:"email" example:"prodemmi@gmail.com"`
	Mobile   string `json:"mobile" example:"+989130243519"`
	Password string `json:"password" binding:"required" example:"secure"`
}

type SignInRes struct {
	UserId  string `json:"user_id" example:"4462ebfc-5f91-4ef0-9cfb-ac6e7687a66e"`
	OtpHash string `json:"otp_hash" example:"7a8c54bf6913a65e19c6cf24d5a2bd5c843b031f"`
}

// SignIn
// @Summary SignIn an exist user
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "null"
// @Param payload body SignInBody true "payload"
// @Success 200 {object} swagger.SuccessJSONResult{data=SignInRes}
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/auth/sign-up [post]
func (c AuthController) SignIn(context *gin.Context) {
	var payload SignInBody

	if err := context.ShouldBindJSON(&payload); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	appOtpType, _ := context.Get("AppOtpType")
	appOtpType = fmt.Sprintf("%v", appOtpType)

	if appOtpType == "email" && payload.Email == "" {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "please enter the email",
		})
		return
	}

	if appOtpType == "mobile" && payload.Mobile == "" {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "please enter the mobile number",
		})
		return
	}

	var user models.User
	var errSignInE error
	var errSignInM error
	if appOtpType == "email" {
		user, errSignInE = c.UserService.GetUserWithEmail(payload.Email)
		if errSignInE != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  errSignInE.Error(),
			})
			return
		}

	} else if appOtpType == "mobile" {
		user, errSignInM = c.UserService.GetUserWithMobile(payload.Mobile)
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errSignInM.Error(),
		})
		return
	}

	if !utils.CheckPasswordHash(payload.Password, user.Password) {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "password is wrong",
		})
		return
	}

	otpHash := utils.CreateAuthHash(user, context.Request.UserAgent())
	_, errOtp := c.OtpService.CreateOTP(user.ID, user.AppID, fmt.Sprintf("%v", appOtpType), otpHash)
	if errOtp != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errOtp.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user_id":  user.ID.String(),
			"otp_hash": otpHash,
		},
	})
}

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
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "user not authorized.",
		})
	}

}

func (c AuthController) OAuth2(context *gin.Context) {}

func (c AuthController) OAuth2Fallback(context *gin.Context) {}

type VerifyOTPBody struct {
	UserID  string `json:"user_id" binding:"required" example:"21636369-8b52-4b4a-97b7-50923ceb3ffd"`
	OtpCode string `json:"otp_code" binding:"required" example:"12345"`
	OtpType string `json:"otp_type" binding:"required" example:"email"`
	OtpHash string `json:"otp_hash" binding:"required" example:"fabecfb8ad5895a6a6309e5a7d1b5c48601d1f9e"`
}

type VerifyOTPRes struct {
	ApiToken   string `json:"api_token" example:"eyJhbGciOiJ..."`
	RedirectTo string `json:"redirect_to" example:"https://google.com"`
}

// VerifyOtp
// @Summary Verify otp
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "null"
// @Param payload body VerifyOTPBody true "payload"
// @Success 200 {object} swagger.SuccessJSONResult{data=VerifyOTPRes}
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/auth/verify-otp [post]
func (c AuthController) VerifyOtp(context *gin.Context) {
	var payload VerifyOTPBody

	if err := context.ShouldBindJSON(&payload); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	user, _ := c.UserService.GetUser(payload.UserID)

	if utils.VerifyAuthHash(payload.OtpHash, *user, context.Request.UserAgent()) {
		verify, err := c.OtpService.CheckOTP(user.ID, user.AppID, payload.OtpCode, payload.OtpType, payload.OtpHash)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		if verify {
			appUrl, _ := context.Get("AppAuthUrl")
			token, errToken := c.UserService.CreateUserToken(user.ID.String())
			if errToken != nil {
				context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"error":  errToken.Error(),
				})
				return
			}
			context.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status": "success",
				"data": gin.H{
					"api_token":   token,
					"redirect_to": appUrl,
				},
			})
			return
		}
		return
	}

	context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"status": "error",
		"error":  "cannot verify otp hash",
	})
}

// Logout
// @Summary Logout user
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "null"
// @Param API-TOKEN header string true "null"
// @Header 308 {string} Location "https://google.com"
// @Success 308
// @failure 500 {object} swagger.ErrorJSONResult{}
// @failure 400 {object} swagger.ErrorJSONResult{}
// @Router /v1/auth/logout [delete]
func (c AuthController) Logout(context *gin.Context) {
	appID, _ := context.Get("AppID")
	userID, _ := context.Get("UserID")

	result, err := c.UserService.DeleteToken(utils.AnyToString(appID), utils.AnyToString(userID))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if result {
		appUrl, _ := context.Get("AppUrl")
		context.JSON(http.StatusPermanentRedirect, utils.AnyToString(appUrl))
		return
	}

	context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"status": "error",
		"error":  "cannot logout",
	})
}

func NewAuthController(authService services.IAuthService, userService services.IUserService, otpService services.IOTPService) IAuthController {
	return AuthController{
		AuthService: authService,
		UserService: userService,
		OtpService:  otpService,
	}
}
