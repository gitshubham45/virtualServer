package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitshubham45/virtualServer/internal/routers"

)

func main(){
	fmt.Println("virtual server in golang")

	router := gin.Default()

	router.GET("/ping" , func (c *gin.Context)  {
		c.JSON(http.StatusOK , gin.H{
			"message" : "pong",
		})
	})

	api := router.Group("/api")

	routers.ServerRouter(api)

	router.Run(":8080")
}