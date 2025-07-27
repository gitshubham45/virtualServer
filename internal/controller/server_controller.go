package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitshubham45/virtualServer/internal/db"
	"github.com/gitshubham45/virtualServer/internal/logger"
	"github.com/gitshubham45/virtualServer/internal/models"
	"github.com/gitshubham45/virtualServer/internal/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var billingRate = map[string]float64{
	"basic": 5.0,
	"plus":  8.0,
	"prime": 12.0,
}

func CreateServer(c *gin.Context) {
	var req struct {
		Region string `json:"region"`
		Type   string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error decoding req : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	newUUID := uuid.New().String()

	var newServer = &models.Server{
		ID:          newUUID,
		BillingRate: float64(billingRate[req.Type]),
		Status:      "running",
		Region:      req.Region,
		Type:        req.Type,
	}

	result := db.DB.Create(&newServer)
	if result == nil {
		log.Println("Error creating server")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creeating server"})
		return
	}

	logger.LogServerEvent(newServer.ID, "SERVER_CREATED", "New server created.", nil, logger.StringPtr(newServer.Status))

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"id":      newServer.ID,
		"status":  newServer.Status,
	})

}

func GetServersData(c *gin.Context) {
	fmt.Println("inside get server")
	serverId := c.Param("id")

	var server models.Server

	result := db.DB.First(&server, "id = ?", serverId)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Server with ID %s not found. \n", serverId)

			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Server with ID '%s' not found.", serverId),
			})
			return
		}

		logger.LogServerEvent(serverId, "SERVER_NOT_FOUND", "server not found.", nil, nil)
		log.Printf("Error fetching server details foe ID '%s' : '%v' \n", serverId, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching server details",
			"error":   result.Error.Error(),
		})
		return
	}

	logger.LogServerEvent(serverId, "SERVER_FOUND", "server found.", logger.StringPtr(server.Status), logger.StringPtr(server.Status))
	c.JSON(http.StatusOK, gin.H{
		"message": "Server details fetched successfully",
		"server":  server,
	})
}

func CompleteAction(c *gin.Context) {
	fmt.Println("Inside CompleteAction handler")

	serverId := c.Param("id")
	log.Printf("Attempting action on server with ID: %s\n", serverId)

	var req struct {
		Action string `json:"action"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body format",
			"error":   err.Error(),
		})
		return
	}

	action := req.Action

	var server models.Server
	result := db.DB.First(&server, "id = ?", serverId)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Server with ID '%s' not found.\n", serverId)
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Server with ID '%s' not found.", serverId),
			})
			return
		}
		log.Printf("Error fetching server details for ID '%s': %v\n", serverId, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching server details",
			"error":   result.Error.Error(),
		})
		return
	}

	originalStatus := server.Status
	newStatus, errorMessage := service.HandleAction(action, originalStatus)

	if errorMessage != "" {
		logger.LogServerEvent(server.ID, "ACTION_DENIED", errorMessage, logger.StringPtr(originalStatus), nil)
		log.Printf("Invalid state transition for server '%s': %s (Current: %s, Action: %s)\n",
			server.ID, errorMessage, originalStatus, action)
		c.JSON(http.StatusConflict, gin.H{
			"message": errorMessage,
		})
		return
	}

	if newStatus != "" {
		server.Status = newStatus
		updateResult := db.DB.Save(&server)
		if updateResult.Error != nil {
			log.Printf("Error saving new status for server '%s': %v\n", server.ID, updateResult.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update server status",
				"error":   updateResult.Error.Error(),
			})
			return
		}

		logger.LogServerEvent(server.ID, "STATUS_CHANGE", fmt.Sprintf("Status changed to '%s'.", newStatus), logger.StringPtr(originalStatus), logger.StringPtr(newStatus))
		
		log.Printf("Server '%s' status changed from '%s' to '%s' via action '%s'.\n",
			server.ID, originalStatus, newStatus, action)
		c.JSON(http.StatusOK, gin.H{
			"message": "Server action completed successfully",
			"server": gin.H{
				"id":        server.ID,
				"status":    server.Status,
				"stoppedAt": server.UpdatedAt,
			},
		})
		return
	}

	logger.LogServerEvent(server.ID, fmt.Sprintf("ACTION_%s_NO_CHANGE", action), fmt.Sprintf("Action '%s' processed, status remains '%s'.", action, originalStatus), logger.StringPtr(originalStatus), nil)
	log.Printf("Action '%s' on server '%s' completed without state change (current status: %s).\n",
		action, server.ID, originalStatus)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Action '%s' processed for server. Status remains '%s'.", action, originalStatus),
		"server": gin.H{
			"id":        server.ID,
			"status":    server.Status,
			"stoppedAt": server.UpdatedAt,
		},
	})
}

func ListServers(c *gin.Context) {
	var servers []models.Server
	result := db.DB.Find(&servers)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Servers not found. \n")
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Servers not found.",
			})
			return
		}

		log.Printf("Error fetching server details : '%v' \n", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching server details",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Server list fetched successfully",
		"server":  servers,
	})
}

func GetLogs(c *gin.Context) {

}
