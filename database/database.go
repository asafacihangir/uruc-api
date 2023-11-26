package database

import (
	"fmt"
	"github.com/org_phoenix/orbey/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

// DB is a global variable to hold the connection to the database.
var DB *gorm.DB

// ConnectDatabase initializes the database connection.
func ConnectDatabase() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı: ", err)
	}

	// User modelini veritabanına map'le
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Veritabanı migrasyonunda hata: ", err)
	}

	fmt.Println("Veritabanına başarıyla bağlanıldı ve User modeli map'lendi.")

}
