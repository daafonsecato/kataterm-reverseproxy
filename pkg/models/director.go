package models

func (store *SessionStore) getMachineHost(machineID int) (string, error) {
	var host string
	query := `SELECT domain FROM machines WHERE id = $1`
	err := store.db.QueryRow(query, machineID).Scan(&host)
	if err != nil {
		return "", err
	}
	return host, nil
}
