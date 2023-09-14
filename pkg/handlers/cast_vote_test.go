package handlers_test

import (
	"bytes"
	"encoding/json"
	goframework "github.com/a-novel/go-framework"
	"github.com/a-novel/votes-service/pkg/handlers"
	"github.com/a-novel/votes-service/pkg/models"
	servicesmocks "github.com/a-novel/votes-service/pkg/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCastVoteHandler(t *testing.T) {
	data := []struct {
		name string

		authorization string

		body interface{}

		shouldCallService     bool
		shouldCallServiceWith models.VoteForm
		serviceResp           *models.VotesSummary
		serviceErr            error

		expect       interface{}
		expectStatus int
	}{
		{
			name:          "Success",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"targetID": goframework.NumberUUID(1).String(),
				"target":   "target",
				"vote":     "up",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.VoteForm{
				TargetID: goframework.NumberUUID(1),
				Target:   "target",
				Vote:     lo.ToPtr(models.VoteValueUp),
			},
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
			name:          "Success/NoVote",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"targetID": goframework.NumberUUID(1).String(),
				"target":   "target",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.VoteForm{
				TargetID: goframework.NumberUUID(1),
				Target:   "target",
			},
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
			name:          "Error/ErrInvalidCredentials",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"targetID": goframework.NumberUUID(1).String(),
				"target":   "target",
				"vote":     "up",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.VoteForm{
				TargetID: goframework.NumberUUID(1),
				Target:   "target",
				Vote:     lo.ToPtr(models.VoteValueUp),
			},
			serviceErr:   goframework.ErrInvalidCredentials,
			expectStatus: http.StatusForbidden,
		},
		{
			name:          "Error/ErrInvalidEntity",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"targetID": goframework.NumberUUID(1).String(),
				"target":   "target",
				"vote":     "up",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.VoteForm{
				TargetID: goframework.NumberUUID(1),
				Target:   "target",
				Vote:     lo.ToPtr(models.VoteValueUp),
			},
			serviceErr:   goframework.ErrInvalidEntity,
			expectStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewCastVoteService(t)

			mrshBody, err := json.Marshal(d.body)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(mrshBody))
			c.Request.Header.Set("Authorization", d.authorization)

			if d.shouldCallService {
				service.
					On("Cast", c, d.authorization, d.shouldCallServiceWith, mock.Anything, mock.Anything).
					Return(d.serviceResp, d.serviceErr)
			}

			handler := handlers.NewCastVoteHandler(service)
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
