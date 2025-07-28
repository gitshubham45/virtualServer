package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gitshubham45/virtualServer/internal/db"
	"github.com/gitshubham45/virtualServer/internal/routers"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("virtual server in golang")


	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .enf file")
	}

	port := os.Getenv("APP_PORT")

	db.InitDB()
	defer db.CloseDB()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api := router.Group("/api")

	routers.ServerRouter(api)

	router.Run(":" + port)
}
