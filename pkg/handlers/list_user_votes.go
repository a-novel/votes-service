package handlers

import (
	"github.com/a-novel/go-apis"
	goframework "github.com/a-novel/go-framework"
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
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
			{goframework.ErrInvalidEntity, http.StatusUnprocessableEntity},
		}, false)
		return
	}

	c.JSON(http.StatusOK, gin.H{"votes": votes})
}
