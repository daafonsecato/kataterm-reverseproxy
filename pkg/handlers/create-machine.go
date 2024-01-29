package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func generateSessionID() string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		// Handle the error according to your application's needs
		// For example, you might want to log the error and return a fallback value
		return "fallback-session-id"
	}
	return newUUID.String()
}

func (controller *SessionController) createMachineHandler(w http.ResponseWriter, r *http.Request) {

	result, err := controller.AWSService.CreateInstance("ami-0c4c061339e1a2038", "t2.micro", "subnet-0f2a559b782561e6d")
	if err != nil {
		fmt.Fprint(w, "Failed to create EC2 machine")
		return
	}
	// Assume you have a way to generate or obtain a session ID
	sessionID := generateSessionID() // Implement this function

	if result != nil {
		instance := result
		instanceID := *instance.InstanceId
		var ipAddress string
		if instance.PrivateIpAddress != nil {
			ipAddress = *instance.PrivateIpAddress
		}

		// Store the machine and session information in the database
		controller.sessionStore.storeMachineAndSession(controller.sessionStore.db, instanceID, ipAddress, sessionID)

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
