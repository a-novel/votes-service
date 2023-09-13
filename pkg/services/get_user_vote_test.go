package services_test

import (
	"context"
	authmocks "github.com/a-novel/auth-service/framework/mocks"
	authmodels "github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/postgresql"
	"github.com/a-novel/go-framework/test"
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

		authClientResp *authmodels.UserTokenStatus
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
			targetID: test.NumberUUID(1),
			target:   "target",
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			shouldCallDAO: true,
			daoResp: &dao.VoteModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(10), baseTime, nil),
				Vote:     models.VoteValueUp,
				UserID:   test.NumberUUID(100),
				TargetID: test.NumberUUID(1),
				Target:   "target",
			},
			expect: &models.Vote{
				ID:        test.NumberUUID(10),
				UpdatedAt: baseTime,
				Vote:      models.VoteValueUp,
				UserID:    test.NumberUUID(100),
				TargetID:  test.NumberUUID(1),
				Target:    "target",
			},
		},
		{
			name:     "Success/Updated",
			tokenRaw: "token",
			targetID: test.NumberUUID(1),
			target:   "target",
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			shouldCallDAO: true,
			daoResp: &dao.VoteModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(10), baseTime, &updateTime),
				Vote:     models.VoteValueUp,
				UserID:   test.NumberUUID(100),
				TargetID: test.NumberUUID(1),
				Target:   "target",
			},
			expect: &models.Vote{
				ID:        test.NumberUUID(10),
				UpdatedAt: updateTime,
				Vote:      models.VoteValueUp,
				UserID:    test.NumberUUID(100),
				TargetID:  test.NumberUUID(1),
				Target:    "target",
			},
		},
		{
			name:     "Error/DAOFailure",
			tokenRaw: "token",
			targetID: test.NumberUUID(1),
			target:   "target",
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			shouldCallDAO: true,
			daoErr:        fooErr,
			expectErr:     fooErr,
		},
		{
			name:           "Error/NotAuthenticated",
			tokenRaw:       "token",
			targetID:       test.NumberUUID(1),
			target:         "target",
			authClientResp: &authmodels.UserTokenStatus{},
			expectErr:      errors.ErrInvalidCredentials,
		},
		{
			name:          "Error/AuthClientFailure",
			tokenRaw:      "token",
			targetID:      test.NumberUUID(1),
			target:        "target",
			authClientErr: fooErr,
			expectErr:     fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			repository := daomocks.NewVotesRepository(t)
			authClient := authmocks.NewClient(t)

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
