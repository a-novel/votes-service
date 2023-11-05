package services_test

import (
	"context"
	"github.com/a-novel/bunovel"
	apiclients "github.com/a-novel/go-apis/clients"
	apiclientsmocks "github.com/a-novel/go-apis/clients/mocks"
	goframework "github.com/a-novel/go-framework"
	"github.com/a-novel/votes-service/pkg/dao"
	daomocks "github.com/a-novel/votes-service/pkg/dao/mocks"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/a-novel/votes-service/pkg/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUserVoteService(t *testing.T) {
	data := []struct {
		name string

		tokenRaw string
		targetID uuid.UUID
		target   string

		authClientResp *apiclients.UserTokenStatus
		authClientErr  error

		shouldCallDAO bool
		daoResp       *dao.VoteModel
		daoErr        error

		expect    *models.Vote
		expectErr error
	}{
		{
			name:     "Success",
			tokenRaw: "token",
			targetID: goframework.NumberUUID(1),
			target:   "target",
			authClientResp: &apiclients.UserTokenStatus{
				OK: true,
				Token: &apiclients.UserToken{
					Payload: apiclients.UserTokenPayload{ID: goframework.NumberUUID(100)},
				},
			},
			shouldCallDAO: true,
			daoResp: &dao.VoteModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(10), baseTime, nil),
				Vote:     models.VoteValueUp,
				UserID:   goframework.NumberUUID(100),
				TargetID: goframework.NumberUUID(1),
				Target:   "target",
			},
			expect: &models.Vote{
				ID:        goframework.NumberUUID(10),
				UpdatedAt: baseTime,
				Vote:      models.VoteValueUp,
				UserID:    goframework.NumberUUID(100),
				TargetID:  goframework.NumberUUID(1),
				Target:    "target",
			},
		},
		{
			name:     "Success/Updated",
			tokenRaw: "token",
			targetID: goframework.NumberUUID(1),
			target:   "target",
			authClientResp: &apiclients.UserTokenStatus{
				OK: true,
				Token: &apiclients.UserToken{
					Payload: apiclients.UserTokenPayload{ID: goframework.NumberUUID(100)},
				},
			},
			shouldCallDAO: true,
			daoResp: &dao.VoteModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(10), baseTime, &updateTime),
				Vote:     models.VoteValueUp,
				UserID:   goframework.NumberUUID(100),
				TargetID: goframework.NumberUUID(1),
				Target:   "target",
			},
			expect: &models.Vote{
				ID:        goframework.NumberUUID(10),
				UpdatedAt: updateTime,
				Vote:      models.VoteValueUp,
				UserID:    goframework.NumberUUID(100),
				TargetID:  goframework.NumberUUID(1),
				Target:    "target",
			},
		},
		{
			name:     "Error/DAOFailure",
			tokenRaw: "token",
			targetID: goframework.NumberUUID(1),
			target:   "target",
			authClientResp: &apiclients.UserTokenStatus{
				OK: true,
				Token: &apiclients.UserToken{
					Payload: apiclients.UserTokenPayload{ID: goframework.NumberUUID(100)},
				},
			},
			shouldCallDAO: true,
			daoErr:        fooErr,
			expectErr:     fooErr,
		},
		{
			name:           "Error/NotAuthenticated",
			tokenRaw:       "token",
			targetID:       goframework.NumberUUID(1),
			target:         "target",
			authClientResp: &apiclients.UserTokenStatus{},
			expectErr:      goframework.ErrInvalidCredentials,
		},
		{
			name:          "Error/AuthClientFailure",
			tokenRaw:      "token",
			targetID:      goframework.NumberUUID(1),
			target:        "target",
			authClientErr: fooErr,
			expectErr:     fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			repository := daomocks.NewVotesRepository(t)
			authClient := apiclientsmocks.NewAuthClient(t)

			authClient.On("IntrospectToken", context.Background(), d.tokenRaw).Return(d.authClientResp, d.authClientErr)

			if d.shouldCallDAO {
				repository.
					On("Get", context.Background(), d.authClientResp.Token.Payload.ID, d.targetID, d.target).
					Return(d.daoResp, d.daoErr)
			}

			service := services.NewGetUserVoteService(repository, authClient)

			resp, err := service.Get(context.Background(), d.tokenRaw, d.targetID, d.target)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, resp)

			repository.AssertExpectations(t)
			authClient.AssertExpectations(t)
		})
	}
}
