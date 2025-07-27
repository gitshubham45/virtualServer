package service

import (
	"fmt"
	"log"
)

const (
	StatusPending    = "pending"
	StatusRunning    = "running"
	StatusStopped    = "stopped"
	StatusTerminated = "terminated"
)

const (
	ActionStart     = "start"
	ActionStop      = "stop"
	ActionReboot    = "reboot"
	ActionTerminate = "terminate"
)

func HandleAction(action string, originalStatus string) (string, string) {
	var newStatus string
	var errorMessage string

	switch action {
	case ActionStart:
		switch originalStatus {
		case StatusStopped:
			newStatus = StatusRunning
		case StatusRunning:
			errorMessage = "Server is already running."
		case StatusTerminated:
			errorMessage = "Cannot start a terminated server."
		case StatusPending:
			errorMessage = "Server is in pending state and cannot be started."
		default:
			errorMessage = fmt.Sprintf("Cannot start server from '%s' status.", originalStatus)
		}

	case ActionStop:
		switch originalStatus {
		case StatusRunning:
			newStatus = StatusStopped
		case StatusStopped:
			errorMessage = "Server is already stopped."
		case StatusTerminated:
			errorMessage = "Cannot stop a terminated server."
		case StatusPending:
			errorMessage = "Server is in pending state and cannot be stopped."
		default:
			errorMessage = fmt.Sprintf("Cannot stop server from '%s' status.", originalStatus)
		}

	case ActionReboot:
		switch originalStatus {
		case StatusRunning:
			newStatus = StatusRunning
			log.Printf("Server is being rebooted. This is often a transient operation.")
		case StatusStopped:
			errorMessage = "Cannot reboot a stopped server. Start it first."
		case StatusTerminated:
			errorMessage = "Cannot reboot a terminated server."
		case StatusPending:
			errorMessage = "Server is in pending state and cannot be rebooted."
		default:
			errorMessage = fmt.Sprintf("Cannot reboot server from '%s' status.", originalStatus)
		}

	case ActionTerminate:
		switch originalStatus {
		case StatusRunning, StatusStopped, StatusPending:
			newStatus = StatusTerminated
		case StatusTerminated:
			errorMessage = "Server is already terminated."
		default:
			errorMessage = fmt.Sprintf("Cannot terminate server from '%s' status.", originalStatus)
		}

	default:
		errorMessage = fmt.Sprintf("Action '%s' is not supported.", action)
		newStatus = ""
	}

	return newStatus, errorMessage
}