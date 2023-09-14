package handlers_test

import (
	"encoding/json"
	goframework "github.com/a-novel/go-framework"
	"github.com/a-novel/votes-service/pkg/handlers"
	"github.com/a-novel/votes-service/pkg/models"
	servicesmocks "github.com/a-novel/votes-service/pkg/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListUserVoteHandler(t *testing.T) {
	data := []struct {
		name string

		authorization string

		query string

		shouldCallService     bool
		shouldCallServiceWith *models.ListUserVotesQuery
		serviceResp           []*models.Vote
		serviceErr            error

		expect       interface{}
		expectStatus int
	}{
		{
			name:              "Success",
			authorization:     "Bearer my-token",
			query:             "?target=target&limit=10&offset=5",
			shouldCallService: true,
			shouldCallServiceWith: &models.ListUserVotesQuery{
				Target: "target",
				Limit:  10,
				Offset: 5,
			},
			serviceResp: []*models.Vote{
				{
					ID:        goframework.NumberUUID(10),
					UpdatedAt: baseTime,
					Vote:      models.VoteValueUp,
					UserID:    goframework.NumberUUID(100),
					TargetID:  goframework.NumberUUID(1),
					Target:    "target",
				},
				{
					ID:        goframework.NumberUUID(20),
					UpdatedAt: updateTime,
					Vote:      models.VoteValueDown,
					UserID:    goframework.NumberUUID(100),
					TargetID:  goframework.NumberUUID(2),
					Target:    "target",
				},
			},
			expect: map[string]interface{}{
				"votes": []interface{}{
					map[string]interface{}{
						"id":        goframework.NumberUUID(10).String(),
						"updatedAt": baseTime.Format(time.RFC3339),
						"vote":      "up",
						"userID":    goframework.NumberUUID(100).String(),
						"targetID":  goframework.NumberUUID(1).String(),
						"target":    "target",
					},
					map[string]interface{}{
						"id":        goframework.NumberUUID(20).String(),
						"updatedAt": updateTime.Format(time.RFC3339),
						"vote":      "down",
						"userID":    goframework.NumberUUID(100).String(),
						"targetID":  goframework.NumberUUID(2).String(),
						"target":    "target",
					},
				},
			},
			expectStatus: http.StatusOK,
		},
		{
			name:              "Error/ErrInvalidCredentials",
			authorization:     "Bearer my-token",
			query:             "?target=target&limit=10&offset=5",
			shouldCallService: true,
			shouldCallServiceWith: &models.ListUserVotesQuery{
				Target: "target",
				Limit:  10,
				Offset: 5,
			},
			serviceErr:   goframework.ErrInvalidCredentials,
			expectStatus: http.StatusForbidden,
		},
		{
			name:              "Error/ErrInvalidEntity",
			authorization:     "Bearer my-token",
			query:             "?target=target&limit=10&offset=5",
			shouldCallService: true,
			shouldCallServiceWith: &models.ListUserVotesQuery{
				Target: "target",
				Limit:  10,
				Offset: 5,
			},
			serviceErr:   goframework.ErrInvalidEntity,
			expectStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewListUserVotesService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/"+d.query, nil)
			c.Request.Header.Set("Authorization", d.authorization)

			if d.shouldCallService {
				service.
					On("List", c, d.authorization, d.shouldCallServiceWith).
					Return(d.serviceResp, d.serviceErr)
			}

			handler := handlers.NewListUserVotesHandler(service)
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
