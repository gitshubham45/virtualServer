package db

import (
	"log"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}