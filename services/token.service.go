package services

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"account/models"
	"gorm.io/gorm"
	"time"
)

type (
	TokenService struct {
		App *gin.Engine
		DB  *gorm.DB
	}
	tokenService interface {
		Create(appID uuid.UUID, userID uuid.UUID) (*models.Token, error)
		FindOne(token string) (*models.Token, error)
	}
)

func (a TokenService) Create(appID uuid.UUID, userID uuid.UUID) (*models.Token, error) {
	var result models.Token

	token := jwt.New(jwt.SigningMethodEdDSA)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(1 * time.Minute)
	claims["app_id"] = appID
	claims["user_id"] = userID

	tokenString, err := token.SignedString("sampleSecretKey")
	if err != nil {
		return nil, err
	}

	errCreate := a.DB.Model(&result).Create(&models.Token{
		AppID:  appID,
		UserID: userID,
		Token:  tokenString,
	}).Error

	if errCreate != nil {
		return nil, errCreate
	}

	return &result, nil
}

func (a TokenService) FindOne(token string) (*models.Token, error) {
	var result models.Token

	err := a.DB.Where("token = ?", token).First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func CreateNewTokenService(app *gin.Engine, db *gorm.DB) *TokenService {
	return &TokenService{
		App: app,
		DB:  db,
	}
}
