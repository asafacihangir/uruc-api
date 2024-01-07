package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/org_phoenix/orbey/database"
	"github.com/org_phoenix/orbey/entity"
	"github.com/org_phoenix/orbey/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// CreateUser endpoint'i için handler
func CreateUser(c *gin.Context) {
	// Kullanıcıdan alınacak verileri tutacak yapı
	var input struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=7,max=15"`
	}

	// JSON verilerini çözümle
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	// Giriş verilerini doğrula
	if err := validate.Struct(input); err != nil {
		// Validator hatalarını düzenli bir formata dönüştür
		var fieldErrors []model.FieldError
		for _, err := range err.(validator.ValidationErrors) {
			fieldError := model.FieldError{
				Field:   err.Field(),
				Message: err.Error(),
			}
			fieldErrors = append(fieldErrors, fieldError)
		}
		c.JSON(http.StatusBadRequest, gin.H{"fieldErrors": fieldErrors})
		return
	}

	// E-posta adresinin tekilliğini kontrol et
	var existingUser entity.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bu e-posta adresiyle zaten bir kullanıcı var."})
		return
	}

	// Şifreyi hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Şifre hashlenirken bir hata oluştu."})
		return
	}

	// User modelini oluştur
	user := entity.User{
		ID:        uuid.New(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashedPassword),
	}

	// Kullanıcıyı veritabanına kaydet
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı kaydedilirken bir hata oluştu."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kullanıcı başarıyla kaydedildi."})
}
