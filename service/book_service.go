package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/org_phoenix/orbey/database"
	"github.com/org_phoenix/orbey/entity"
	"github.com/org_phoenix/orbey/model"
	"gorm.io/gorm"
	"net/http"
)

func CreateBook(c *gin.Context) {
	// Kullanıcıdan alınacak verileri tutacak yapı
	var input struct {
		Title           string `json:"title" validate:"required"`
		Author          string `json:"author" validate:"required"`
		Publisher       string `json:"publisher" validate:"required"`
		PublicationDate string `json:"publicationDate" validate:"required"`
		Isbn            string `json:"isbn" validate:"required"`
		Genre           string `json:"genre" validate:"required"`
	}

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
	var existingBook entity.Book
	if err := database.DB.Where("isbn = ?", input.Isbn).First(&existingBook).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bu isbn numarasıyla zaten bir kitap var."})
		return
	}

	book := entity.Book{
		ID:              uuid.New(),
		Title:           input.Title,
		Author:          input.Author,
		Publisher:       input.Publisher,
		PublicationDate: input.PublicationDate,
		Isbn:            input.Isbn,
		Genre:           input.Genre,
	}

	// Kullanıcıyı veritabanına kaydet
	if err := database.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitap kaydedilirken bir hata oluştu."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kitap başarıyla kaydedildi."})
}

func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Title           string `json:"title" validate:"required"`
		Author          string `json:"author" validate:"required"`
		Publisher       string `json:"publisher" validate:"required"`
		PublicationDate string `json:"publicationDate" validate:"required"`
		Isbn            string `json:"isbn" validate:"required"`
		Genre           string `json:"genre" validate:"required"`
	}

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

	// Veritabanında kitabı bul
	var book entity.Book
	if err := database.DB.First(&book, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Kitabı güncelle
	if err := database.DB.Model(&book).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitap güncellenirken bir hata oluştu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kitap başarıyla güncellendi."})
}

func ListBooks(c *gin.Context) {
	var books []entity.Book
	if err := database.DB.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, books)
}

func FindBookById(c *gin.Context) {
	id := c.Param("id")

	var book entity.Book
	if err := database.DB.First(&book, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, book)
}

func DeleteBookById(c *gin.Context) {
	id := c.Param("id")

	// Önce UUID ile kitabı bul
	var book entity.Book
	if err := database.DB.Where("id = ?", id).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kitap bulunamadı"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitap sorgulanırken bir hata oluştu"})
		}
		return
	}

	// Kitabı sil
	if err := database.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitap silinirken bir hata oluştu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kitap başarıyla silindi."})
}
