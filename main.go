package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/org_phoenix/orbey/database"
	"github.com/org_phoenix/orbey/service"
	"log"
)

func init() {
	// .env dosyasını yükler
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	// Veritabanı bağlantısını başlat
	database.ConnectDatabase()

	// Gin router'ını başlat
	router := gin.Default()

	// Kullanıcı kaydetme route'u
	router.POST("/user", service.CreateUser)

	// GET yolu için bir handler (işleyici) tanımla
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	// Servisi 8080 portunda çalıştır
	router.Run(":8080")
}
