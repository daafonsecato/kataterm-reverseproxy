package handlers

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(database *sql.DB) *SessionStore {
	database.InitDB()
	db, err := database.InitDB()
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to open database: %v", err)
		panic(errorMsg)
	}

	err = db.Ping()
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to ping database: %v", err)
		panic(errorMsg)
	}

	if db == nil {
		panic("db is nil")
	}
	return &SessionStore{
		db: db,
	}
}

func (store *SessionStore) Close() {
	store.db.Close()
}

func storeMachineAndSession(db *sql.DB, awsInstanceID, ipAddress string) error {
	// Start a transaction
	tx, err := db.Begin()
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
