package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/org_phoenix/orbey/database"
	"github.com/org_phoenix/orbey/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken         string `json:"accessToken"`
	ExpirationInSeconds int64  `json:"expirationInSeconds"`
	ExpirationDate      string `json:"expirationDate"`
}

func Login(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtSecret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sunucu yapılandırma hatası"})
		return
	}

	var existingUser entity.User
	if result := database.DB.Where("email = ?", loginReq.Email).First(&existingUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Veritabanı sorgulama hatası"})
		}
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginReq.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username": loginReq.Email,
		"exp":      expirationTime.Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	response := LoginResponse{
		AccessToken:         tokenString,
		ExpirationInSeconds: int64(time.Until(expirationTime).Seconds()),
		ExpirationDate:      expirationTime.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)

}
