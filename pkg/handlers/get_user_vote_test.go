package handlers_test

import (
	"encoding/json"
	"github.com/a-novel/bunovel"
	goframework "github.com/a-novel/go-framework"
	"github.com/a-novel/votes-service/pkg/handlers"
	"github.com/a-novel/votes-service/pkg/models"
	servicesmocks "github.com/a-novel/votes-service/pkg/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserVoteHandler(t *testing.T) {
	data := []struct {
		name string

		authorization string

		query string

		shouldCallService             bool
		shouldCallServiceWithTargetID uuid.UUID
		shouldCallServiceWithTarget   string
		serviceResp                   *models.Vote
		serviceErr                    error

		expect       interface{}
		expectStatus int
	}{
		{
			name:                          "Success",
			authorization:                 "Bearer my-token",
			query:                         "?targetID=01010101-0101-0101-0101-010101010101&target=target",
			shouldCallService:             true,
			shouldCallServiceWithTargetID: goframework.NumberUUID(1),
			shouldCallServiceWithTarget:   "target",
			serviceResp: &models.Vote{
				ID:        goframework.NumberUUID(10),
				UpdatedAt: baseTime,
				Vote:      models.VoteValueUp,
				UserID:    goframework.NumberUUID(100),
				TargetID:  goframework.NumberUUID(1),
				Target:    "target",
			},
			expect: map[string]interface{}{
				"id":        goframework.NumberUUID(10).String(),
				"updatedAt": baseTime.Format(time.RFC3339),
				"vote":      "up",
				"userID":    goframework.NumberUUID(100).String(),
				"targetID":  goframework.NumberUUID(1).String(),
				"target":    "target",
			},
			expectStatus: http.StatusOK,
		},
		{
			name:                          "Error/ErrInvalidCredentials",
			authorization:                 "Bearer my-token",
			query:                         "?targetID=01010101-0101-0101-0101-010101010101&target=target",
			shouldCallService:             true,
			shouldCallServiceWithTargetID: goframework.NumberUUID(1),
			shouldCallServiceWithTarget:   "target",
			serviceErr:                    goframework.ErrInvalidCredentials,
			expectStatus:                  http.StatusForbidden,
		},
		{
			name:                          "Error/ErrNotFound",
			authorization:                 "Bearer my-token",
			query:                         "?targetID=01010101-0101-0101-0101-010101010101&target=target",
			shouldCallService:             true,
			shouldCallServiceWithTargetID: goframework.NumberUUID(1),
			shouldCallServiceWithTarget:   "target",
			serviceErr:                    bunovel.ErrNotFound,
			expectStatus:                  http.StatusNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewGetUserVoteService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/"+d.query, nil)
			c.Request.Header.Set("Authorization", d.authorization)

			if d.shouldCallService {
				service.
					On("Get", c, d.authorization, d.shouldCallServiceWithTargetID, d.shouldCallServiceWithTarget).
					Return(d.serviceResp, d.serviceErr)
			}

			handler := handlers.NewGetUserVoteHandler(service)
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
