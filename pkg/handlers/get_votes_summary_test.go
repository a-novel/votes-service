package handlers_test

import (
	"encoding/json"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/test"
	"github.com/a-novel/votes-service/pkg/handlers"
	"github.com/a-novel/votes-service/pkg/models"
	servicesmocks "github.com/a-novel/votes-service/pkg/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetVotesSummaryHandler(t *testing.T) {
	data := []struct {
		name string

		query string

		shouldCallService             bool
		shouldCallServiceWithTargetID uuid.UUID
		shouldCallServiceWithTarget   string
		serviceResp                   *models.VotesSummary
		serviceErr                    error

		expect       interface{}
		expectStatus int
	}{
		{
			name:                          "Success",
			query:                         "?targetID=01010101-0101-0101-0101-010101010101&target=target",
			shouldCallService:             true,
			shouldCallServiceWithTargetID: test.NumberUUID(1),
			shouldCallServiceWithTarget:   "target",
			serviceResp: &models.VotesSummary{
				UpVotes:   128,
				DownVotes: 64,
			},
			expect: map[string]interface{}{
				"upVotes":   float64(128),
				"downVotes": float64(64),
			},
			expectStatus: http.StatusOK,
		},
		{
			name:                          "Error/ErrNotFound",
			query:                         "?targetID=01010101-0101-0101-0101-010101010101&target=target",
			shouldCallService:             true,
			shouldCallServiceWithTargetID: test.NumberUUID(1),
			shouldCallServiceWithTarget:   "target",
			serviceErr:                    errors.ErrNotFound,
			expectStatus:                  http.StatusNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewGetVotesSummaryService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/"+d.query, nil)

			if d.shouldCallService {
				service.
					On("Get", c, d.shouldCallServiceWithTargetID, d.shouldCallServiceWithTarget).
					Return(d.serviceResp, d.serviceErr)
			}

			handler := handlers.NewGetVotesSummaryHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code, c.Errors.String())
			if d.expect != nil {
				var body interface{}
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
				require.Equal(t, d.expect, body)
			}

			service.AssertExpectations(t)
		})
	}
}
