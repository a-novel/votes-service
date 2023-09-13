package handlers

import (
	auth "github.com/a-novel/auth-service/framework"
	forum "github.com/a-novel/forum-service/framework"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"net/http"
)

type HealthCheckHandler interface {
	Handle(c *gin.Context)
}

func NewHealthCheckHandler(db *bun.DB, authClient auth.Client, forumClient forum.Client) HealthCheckHandler {
	return &healthCheckHandlerImpl{db: db, authClient: authClient}
}

type healthCheckHandlerImpl struct {
	db          *bun.DB
	authClient  auth.Client
	forumClient forum.Client
}

type HealthCheckResponse struct {
	Database struct {
		Available bool   `json:"available"`
		Error     string `json:"error,omitempty"`
	} `json:"database"`
	Clients struct {
		Auth struct {
			Available bool   `json:"available"`
			Error     string `json:"error,omitempty"`
		} `json:"auth"`
		Forum struct {
			Available bool   `json:"available"`
			Error     string `json:"error,omitempty"`
		} `json:"forum"`
	} `json:"clients"`
}

func (h *healthCheckHandlerImpl) Handle(c *gin.Context) {
	res := new(HealthCheckResponse)

	dbErr := h.db.PingContext(c)
	res.Database.Available = dbErr == nil
	if dbErr != nil {
		res.Database.Error = dbErr.Error()
	}

	// Internal api has no auth client.
	if h.authClient != nil {
		authClientErr := h.authClient.Ping()
		res.Clients.Auth.Available = authClientErr == nil
		if authClientErr != nil {
			res.Clients.Auth.Error = authClientErr.Error()
		}
	}

	// External api has no forum client.
	if h.forumClient != nil {
		forumClientErr := h.forumClient.Ping()
		res.Clients.Forum.Available = forumClientErr == nil
		if forumClientErr != nil {
			res.Clients.Forum.Error = forumClientErr.Error()
		}
	}

	c.JSON(http.StatusOK, res)
}
