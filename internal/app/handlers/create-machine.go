package handlers

import (
	"database/sql"
	"github.com/david8128/kataterm-reverseproxy/internal/app/models"
	"github.com/david8128/kataterm-reverseproxy/internal/database"
	"github.com/david8128/kataterm-reverseproxy/internal/app/services"
	"net/http"
)

type SessionStore struct {
	db *sql.DB
}

type SessionController struct {
	sessionStore        *models.SessionStore
}

func NewSessionController() *SessionController {
	db, err := database.InitDB()
	if err != nil {
		panic("Error initializing DB")
	}

	sessionStore := models.NewSessionStore(db)

	return &SessionController{
		sessionStore:        sessionStore,
	}
}
func (controller *SessionController) createMachineHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// Assume you have a way to generate or obtain a session ID
	sessionID := generateSessionID() // Implement this function

	if len(result.Instances) > 0 {
		instance := result.Instances[0]
		instanceID := *instance.InstanceId
		var ipAddress string
		if instance.PrivateIpAddress != nil {
			ipAddress = *instance.PrivateIpAddress
		}

		// Store the machine and session information in the database
		storeMachineAndSession(db, instanceID, ipAddress, sessionID)

		if ipAddress != "" {
			fmt.Fprintf(w, "EC2 machine created with IP address: %s", ipAddress)
			// Generate the domain with session ID as a subdomain
			domain := fmt.Sprintf("%s.terminal.kataterm.com", sessionID)

			// Send the domain in the response body
			fmt.Fprintf(w, "EC2 machine created with domain: %s", domain)

		} else {
			fmt.Fprint(w, "EC2 machine created, but IP address is not available")
		}
		return
	}

	fmt.Fprint(w, "Failed to create EC2 machine")
}
