package handlers

import (
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/a-novel/votes-service/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetUserVoteHandler interface {
	Handle(c *gin.Context)
}

func NewGetUserVoteHandler(service services.GetUserVoteService) GetUserVoteHandler {
	return &getUserVoteHandlerImpl{
		service: service,
	}
}

type getUserVoteHandlerImpl struct {
	service services.GetUserVoteService
}

func (h *getUserVoteHandlerImpl) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	query := new(models.GetUserVoteQuery)
	if err := c.BindQuery(query); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	vote, err := h.service.Get(c, token, query.TargetID.Value(), query.Target)
	if err != nil {
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidCredentials, http.StatusForbidden},
			{errors.ErrNotFound, http.StatusNotFound},
		})
		return
	}

	c.JSON(http.StatusOK, vote)
}
