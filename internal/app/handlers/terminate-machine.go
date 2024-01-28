package handlers

import (
	"database/sql"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/david8128/kataterm-reverseproxy/internal/app/services"
	"net/http"
)

type TerminationRequest struct {
	SessionID string `json:"session_id"`
}

func terminateMachineHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TerminationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	if req.SessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}
	// Use the terminateInstanceBySessionID function to terminate the machine
	err := terminateInstanceBySessionID(db, sessionID)
	if err != nil {
		log.Printf("Failed to terminate machine: %v", err)
		http.Error(w, fmt.Sprintf("Error terminating machine: %v", err), http.StatusInternalServerError)
		return
	}

	// Send a success response
	fmt.Fprintf(w, "Machine terminated successfully for session ID: %s", sessionID)
}
