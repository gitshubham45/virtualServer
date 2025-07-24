package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gitshubham45/virtualServer/internal/controller"
)

func ServerRouter(api *gin.RouterGroup){
	api.POST("/server" , controller.CreateServer)
	api.GET("/server/:id" , controller.GetServersData)
	api.POST("/servers/:id/action" , controller.CompleteAction)
	api.GET("/servers" , controller.ListServers)
	api.GET("/servers/:id/logs" , controller.GetLogs)
}