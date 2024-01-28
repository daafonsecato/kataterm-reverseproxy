
func getAWSInstanceID(db *sql.DB, sessionID string) (string, error) {
    var awsInstanceID string
    query := `SELECT aws_instance_id FROM machines JOIN sessions ON machines.id = sessions.machine_id WHERE sessions.session_id = $1`
    err := db.QueryRow(query, sessionID).Scan(&awsInstanceID)
    if err != nil {
        return "", err
    }
    return awsInstanceID, nil
}