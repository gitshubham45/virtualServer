package db

import (
	"fmt"
	"log"
	"os"

	"github.com/gitshubham45/virtualServer/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	timeZone := os.Getenv("DB_TIME_ZONE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		timeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.AutoMigrate(&models.Server{})
	if err != nil {
		log.Fatalf("Failed to auto migrate schemas : %v" , err)
	}
	fmt.Println("Database schema auto-migrated successfully for server model")

	DB = db
	log.Println("Connected to PostgreSql")
}

func CloseDB() {
	if DB != nil {
		sqlDb, err := DB.DB()
		if err != nil {
			log.Printf("Error getting *sql.DB for closing : %v", err)
		}
		if err := sqlDb.Close(); err != nil {
			log.Printf("Error closing db connection : %v", err)
		} else {
			fmt.Println("Database connection closed")
		}
	}
}
