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
			shouldCallServiceWithTargetID: test.NumberUUID(1),
			shouldCallServiceWithTarget:   "target",
			serviceResp: &models.Vote{
				ID:        test.NumberUUID(10),
				UpdatedAt: baseTime,
				Vote:      models.VoteValueUp,
				UserID:    test.NumberUUID(100),
				TargetID:  test.NumberUUID(1),
				Target:    "target",
			},
			expect: map[string]interface{}{
				"id":        test.NumberUUID(10).String(),
				"updatedAt": baseTime.Format(time.RFC3339),
				"vote":      "up",
				"userID":    test.NumberUUID(100).String(),
				"targetID":  test.NumberUUID(1).String(),
				"target":    "target",
			},
			expectStatus: http.StatusOK,
		},
		{
			name:                          "Error/ErrInvalidCredentials",
			authorization:                 "Bearer my-token",
			query:                         "?targetID=01010101-0101-0101-0101-010101010101&target=target",
			shouldCallService:             true,
			shouldCallServiceWithTargetID: test.NumberUUID(1),
			shouldCallServiceWithTarget:   "target",
			serviceErr:                    errors.ErrInvalidCredentials,
			expectStatus:                  http.StatusForbidden,
		},
		{
			name:                          "Error/ErrNotFound",
			authorization:                 "Bearer my-token",
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
