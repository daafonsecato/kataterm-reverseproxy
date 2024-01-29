package models

import (
	"database/sql"

	"github.com/google/uuid"
)

func (store *SessionStore) storeMachineAndSession(db *sql.DB, awsInstanceID, ipAddress string) error {
	// Start a transaction
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}

	// Insert into machines table
	var machineID int
	insertMachineQuery := `INSERT INTO machines (aws_instance_id, status, domain) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(insertMachineQuery, awsInstanceID, "pending", ipAddress).Scan(&machineID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Generate a new UUID for the session
	sessionID, err := uuid.NewUUID()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert into sessions table
	insertSessionQuery := `INSERT INTO sessions (session_id, machine_id) VALUES ($1, $2)`
	_, err = tx.Exec(insertSessionQuery, sessionID, machineID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit()
}
