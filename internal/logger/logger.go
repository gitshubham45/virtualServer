package logger

import (
	"log"

	"github.com/gitshubham45/virtualServer/internal/db"
	"github.com/gitshubham45/virtualServer/internal/models"
	"github.com/google/uuid"
)

func LogServerEvent(serverID, eventType, message string, oldStatus, newStatus *string) {
	newUUID := uuid.New().String()
	logEntry := models.ServerLog{
		ID:        newUUID,
		ServerID:  serverID,
		EventType: eventType,
		Message:   message,
	}

	if oldStatus != nil {
		logEntry.OldStatus = *oldStatus
	}
	if newStatus != nil {
		logEntry.NewStatus = *newStatus
	}

	if err := db.DB.Create(&logEntry).Error; err != nil {
		log.Printf("WARNING: Failed to save server log for server %s (Event: %s): %v\n", serverID, eventType, err)
	} else {
		log.Printf("Server log saved: ServerID=%s, EventType=%s, Message='%s'\n", serverID, eventType, message)
	}
}

func StringPtr(s string) *string {
	return &s
}
