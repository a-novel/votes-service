package handlers

import (
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/a-novel/votes-service/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListUserVotesHandler interface {
	Handle(c *gin.Context)
}

func NewListUserVotesHandler(service services.ListUserVotesService) ListUserVotesHandler {
	return &listUserVotesHandlerImpl{
		service: service,
	}
}

type listUserVotesHandlerImpl struct {
	service services.ListUserVotesService
}

func (h *listUserVotesHandlerImpl) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	query := new(models.ListUserVotesQuery)
	if err := c.BindQuery(query); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	votes, err := h.service.List(c, token, query)
	if err != nil {
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidCredentials, http.StatusForbidden},
			{errors.ErrInvalidEntity, http.StatusUnprocessableEntity},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"votes": votes})
}
