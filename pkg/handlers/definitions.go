package handlers

import (
	db "github.com/daafonsecato/kataterm-reverseproxy/internal/database"
	"github.com/daafonsecato/kataterm-reverseproxy/pkg/models"
	"github.com/daafonsecato/kataterm-reverseproxy/pkg/services"
)

type SessionController struct {
	sessionStore *models.SessionStore
	AWSService   *services.AWSService
}

func NewSessionController() *SessionController {
	db, err := db.InitDB()
	if err != nil {
		panic("Error initializing DB")
	}

	sessionStore := models.NewSessionStore(db)
	sess := services.NewAWSService()

	return &SessionController{
		sessionStore: sessionStore,
		AWSService:   sess,
	}
}
