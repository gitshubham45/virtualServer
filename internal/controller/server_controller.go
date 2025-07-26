package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitshubham45/virtualServer/internal/db"
	"github.com/gitshubham45/virtualServer/internal/models"
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

		log.Printf("Error fetching server details foe ID '%s' : '%v' \n", serverId, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching server details",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Server details fetched successfully",
		"server":  server,
	})
}

func CompleteAction(c *gin.Context) {
	serverId := c.Param("action")

	var req struct {
		Action string `json:"action"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error decoding req : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	action := req.Action


	// check if the action is valid or not out of 4 option -
	//  "pending," "running," "stopped," "terminated"

	// fetch the server using serverId

	// status - running  - action - stop  - valid
	// status - stopped - action - start - valid
	// status - running - action - reebot - valid
	// status - running - action - terminate - valid
	// status - stopped - action - terminate - valid
	// status - terminated - action - start - Invalid - return server is terminatyed
	// status - terminated - action - stop - Invalid - retunrn servet is terminated
	// status - running - action - start - Invalid - return already running
	// status - stopped - action - stopped - Invalid - return already stopped
	// status - terminated - action - terminate - Invalid - return already terminated

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
