package handlers

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/cssivision/reverseproxy"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"net/url"
)

func getMachineHost(db *sql.DB, machineID int) (string, error) {
	var host string
	query := `SELECT domain FROM machines WHERE id = $1`
	err := db.QueryRow(query, machineID).Scan(&host)
	if err != nil {
		return "", err
	}
	return host, nil
}

func customDirector(db *sql.DB) func(req *http.Request) {
	return func(req *http.Request) {
		host := req.Host
		subdomains := strings.Split(host, ".")

		// Assuming the format is [sessionID].[service].terminal.kataterm.com
		if len(subdomains) < 4 {
			log.Println("Invalid host:", host)
			return // Or set a default targetHost
		}

		sessionID := subdomains[0]
		service := subdomains[1]

		// Get the host for the machine
		targetHost, err := getMachineHost(db, sessionID)
		if err != nil {
			log.Println("Error getting machine host:", err)
			return // Or set a default targetHost
		}

		var port string
		switch service {
		case "backend":
			port = "8000"
		case "ttyd":
			port = "7681"
		case "code-server":
			port = "8080"
		// Add more cases as needed for different services
		default:
			log.Println("Service not recognized:", service)
			return // Or set a default port
		}

		// Modify the request to route to the target host and port
		req.URL.Host = targetHost + ":" + port
		// Other modifications to the request as needed
	}
}
