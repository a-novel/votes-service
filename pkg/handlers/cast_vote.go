package handlers

import (
	"github.com/a-novel/go-apis"
	goframework "github.com/a-novel/go-framework"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/a-novel/votes-service/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type CastVoteHandler interface {
	Handle(c *gin.Context)
}

func NewCastVoteHandler(service services.CastVoteService) CastVoteHandler {
	return &castVoteHandlerImpl{
		service: service,
	}
}

type castVoteHandlerImpl struct {
	service services.CastVoteService
}

func (h *castVoteHandlerImpl) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	request := new(models.VoteForm)
	if err := c.BindJSON(request); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	summary, err := h.service.Cast(c, token, *request, uuid.New(), time.Now())
	if err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
			{goframework.ErrInvalidEntity, http.StatusUnprocessableEntity},
		}, true)
		return
	}

	c.JSON(http.StatusOK, summary)
}
