package handlers

import (
	"fmt"
	"net/http"
)

func (controller *SessionController) CreateMachineHandler(w http.ResponseWriter, r *http.Request) {

	result, err := controller.AWSService.CreateInstance("ami-0c4c061339e1a2038", "t2.medium", "subnet-0f2a559b782561e6d")
	if err != nil {
		fmt.Fprint(w, "Failed to create EC2 machine")
		return
	}

	if result != nil {
		instance := result
		instanceID := *instance.InstanceId
		var ipAddress string
		if instance.PrivateIpAddress != nil {
			ipAddress = *instance.PrivateIpAddress
		}

		// Store the machine and session information in the database
		sessionID, err := controller.sessionStore.StoreMachineAndSession(instanceID, ipAddress)
		if err != nil {
			fmt.Fprint(w, "Failed to create EC2 machine")
			return
		}
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
