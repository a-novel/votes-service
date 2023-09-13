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
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListUserVotesService(t *testing.T) {
	data := []struct {
		name string

		tokenRaw string
		query    *models.ListUserVotesQuery

		authClientResp *authmodels.UserTokenStatus
		authClientErr  error

		shouldCallDAO bool
		daoResp       []*dao.VoteModel
		daoErr        error

		expect    []*models.Vote
		expectErr error
	}{
		{
			name:     "Success",
			tokenRaw: "token",
			query: &models.ListUserVotesQuery{
				Target: "target",
				Limit:  10,
				Offset: 5,
			},
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			shouldCallDAO: true,
			daoResp: []*dao.VoteModel{
				{
					Metadata: postgresql.NewMetadata(test.NumberUUID(10), baseTime, nil),
					Vote:     models.VoteValueUp,
					UserID:   test.NumberUUID(100),
					TargetID: test.NumberUUID(1),
					Target:   "target",
				},
				{
					Metadata: postgresql.NewMetadata(test.NumberUUID(20), baseTime, &updateTime),
					Vote:     models.VoteValueDown,
					UserID:   test.NumberUUID(100),
					TargetID: test.NumberUUID(3),
					Target:   "target",
				},
			},
			expect: []*models.Vote{
				{
					ID:        test.NumberUUID(10),
					UpdatedAt: baseTime,
					Vote:      models.VoteValueUp,
					UserID:    test.NumberUUID(100),
					TargetID:  test.NumberUUID(1),
					Target:    "target",
				},
				{
					ID:        test.NumberUUID(20),
					UpdatedAt: updateTime,
					Vote:      models.VoteValueDown,
					UserID:    test.NumberUUID(100),
					TargetID:  test.NumberUUID(3),
					Target:    "target",
				},
			},
		},
		{
			name:     "Success/NoResults",
			tokenRaw: "token",
			query: &models.ListUserVotesQuery{
				Target: "target",
				Limit:  10,
				Offset: 5,
			},
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			shouldCallDAO: true,
			expect:        []*models.Vote{},
		},
		{
			name:     "Error/DAOFailure",
			tokenRaw: "token",
			query: &models.ListUserVotesQuery{
				Target: "target",
				Limit:  10,
				Offset: 5,
			},
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
			name:     "Error/LimitTooHigh",
			tokenRaw: "token",
			query: &models.ListUserVotesQuery{
				Target: "target",
				Limit:  services.MaxSearchLimit + 1,
				Offset: 5,
			},
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			expectErr: errors.ErrInvalidEntity,
		},
		{
			name:     "Error/NoLimit",
			tokenRaw: "token",
			query: &models.ListUserVotesQuery{
				Target: "target",
				Offset: 5,
			},
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			expectErr: errors.ErrInvalidEntity,
		},
		{
			name:     "Error/NotAuthenticated",
			tokenRaw: "token",
			query: &models.ListUserVotesQuery{
				Target: "target",
				Limit:  10,
				Offset: 5,
			},
			authClientResp: &authmodels.UserTokenStatus{},
			expectErr:      errors.ErrInvalidCredentials,
		},
		{
			name:     "Error/AuthClientFailure",
			tokenRaw: "token",
			query: &models.ListUserVotesQuery{
				Target: "target",
				Limit:  10,
				Offset: 5,
			},
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
					On("ListUserVotes", context.Background(), d.authClientResp.Token.Payload.ID, d.query.Target, d.query.Limit, d.query.Offset).
					Return(d.daoResp, d.daoErr)
			}

			service := services.NewListUserVotesService(repository, authClient)

			resp, err := service.List(context.Background(), d.tokenRaw, d.query)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, resp)

			repository.AssertExpectations(t)
			authClient.AssertExpectations(t)
		})
	}
}
