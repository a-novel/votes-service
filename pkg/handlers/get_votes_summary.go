package handlers

import (
	"github.com/a-novel/bunovel"
	"github.com/a-novel/go-apis"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/a-novel/votes-service/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetVotesSummaryHandler interface {
	Handle(c *gin.Context)
}

func NewGetVotesSummaryHandler(service services.GetVotesSummaryService) GetVotesSummaryHandler {
	return &getVotesSummaryHandlerImpl{
		service: service,
	}
}

type getVotesSummaryHandlerImpl struct {
	service services.GetVotesSummaryService
}

func (h *getVotesSummaryHandlerImpl) Handle(c *gin.Context) {
	query := new(models.GetVotesSummaryQuery)
	if err := c.BindQuery(query); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	summary, err := h.service.Get(c, query.TargetID.Value(), query.Target)
	if err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{bunovel.ErrNotFound, http.StatusNotFound},
		}, false)
		return
	}

	c.JSON(http.StatusOK, summary)
}
