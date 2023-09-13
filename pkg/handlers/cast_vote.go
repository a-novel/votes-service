package handlers

import (
	"github.com/a-novel/go-framework/errors"
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
		if httpErr := errors.AsHTTPClientErr(err); httpErr != nil {
			_ = c.AbortWithError(httpErr.Code, httpErr)
			return
		}

		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidCredentials, http.StatusForbidden},
			{errors.ErrInvalidEntity, http.StatusUnprocessableEntity},
		})
		return
	}

	c.JSON(http.StatusOK, summary)
}
