package handlers

import (
	"github.com/a-novel/bunovel"
	"github.com/a-novel/go-apis"
	goframework "github.com/a-novel/go-framework"
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
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
			{bunovel.ErrNotFound, http.StatusNotFound},
		}, false)
		return
	}

	c.JSON(http.StatusOK, vote)
}
