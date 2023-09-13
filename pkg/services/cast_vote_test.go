package services_test

import (
	"context"
	authmocks "github.com/a-novel/auth-service/framework/mocks"
	authmodels "github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/test"
	"github.com/a-novel/votes-service/pkg/dao"
	daomocks "github.com/a-novel/votes-service/pkg/dao/mocks"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/a-novel/votes-service/pkg/services"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCastVoteService(t *testing.T) {
	data := []struct {
		name string

		tokenRaw string
		form     models.VoteForm
		id       uuid.UUID
		now      time.Time

		authClientResp *authmodels.UserTokenStatus
		authClientErr  error

		clientName string
		clientErr  error

		shouldCallDAO bool

		castErr error

		shouldCallGetSummary bool
		summary              *dao.VotesSummaryModel
		summaryErr           error

		expect    *models.VotesSummary
		expectErr error
	}{
		{
			name:     "Success/UpVote",
			tokenRaw: "token",
			form: models.VoteForm{
				TargetID: test.NumberUUID(1),
				Target:   "target",
				Vote:     lo.ToPtr(models.VoteValueUp),
			},
			id:         test.NumberUUID(10),
			now:        baseTime,
			clientName: "target",
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			shouldCallDAO:        true,
			shouldCallGetSummary: true,
			summary: &dao.VotesSummaryModel{
				TargetID:  test.NumberUUID(1),
				Target:    "target",
				UpVotes:   128,
				DownVotes: 64,
			},
			expect: &models.VotesSummary{
				UpVotes:   128,
				DownVotes: 64,
			},
		},
		{
			name:     "Success/DownVote",
			tokenRaw: "token",
			form: models.VoteForm{
				TargetID: test.NumberUUID(1),
				Target:   "target",
				Vote:     lo.ToPtr(models.VoteValueDown),
			},
			id:         test.NumberUUID(10),
			now:        baseTime,
			clientName: "target",
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			shouldCallDAO:        true,
			shouldCallGetSummary: true,
			summary: &dao.VotesSummaryModel{
				TargetID:  test.NumberUUID(1),
				Target:    "target",
				UpVotes:   128,
				DownVotes: 64,
			},
			expect: &models.VotesSummary{
				UpVotes:   128,
				DownVotes: 64,
			},
		},
		{
			name:     "Success/NoVote",
			tokenRaw: "token",
			form: models.VoteForm{
				TargetID: test.NumberUUID(1),
				Target:   "target",
			},
			id:         test.NumberUUID(10),
			now:        baseTime,
			clientName: "target",
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			shouldCallDAO:        true,
			shouldCallGetSummary: true,
			summary: &dao.VotesSummaryModel{
				TargetID:  test.NumberUUID(1),
				Target:    "target",
				UpVotes:   128,
				DownVotes: 64,
			},
			expect: &models.VotesSummary{
				UpVotes:   128,
				DownVotes: 64,
			},
		},
		{
			name:     "Error/TargetCallFailure",
			tokenRaw: "token",
			form: models.VoteForm{
				TargetID: test.NumberUUID(1),
				Target:   "target",
				Vote:     lo.ToPtr(models.VoteValueUp),
			},
			id:         test.NumberUUID(10),
			now:        baseTime,
			clientName: "target",
			clientErr:  fooErr,
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			shouldCallDAO:        true,
			shouldCallGetSummary: true,
			summary: &dao.VotesSummaryModel{
				TargetID:  test.NumberUUID(1),
				Target:    "target",
				UpVotes:   128,
				DownVotes: 64,
			},
			expectErr: fooErr,
		},
		{
			name:     "Error/GetSummaryFailure",
			tokenRaw: "token",
			form: models.VoteForm{
				TargetID: test.NumberUUID(1),
				Target:   "target",
				Vote:     lo.ToPtr(models.VoteValueUp),
			},
			id:         test.NumberUUID(10),
			now:        baseTime,
			clientName: "target",
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			shouldCallDAO:        true,
			shouldCallGetSummary: true,
			summaryErr:           fooErr,
			expectErr:            fooErr,
		},
		{
			name:     "Error/CastVoteFailure",
			tokenRaw: "token",
			form: models.VoteForm{
				TargetID: test.NumberUUID(1),
				Target:   "target",
				Vote:     lo.ToPtr(models.VoteValueUp),
			},
			id:         test.NumberUUID(10),
			now:        baseTime,
			clientName: "target",
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			shouldCallDAO: true,
			castErr:       fooErr,
			expectErr:     fooErr,
		},
		{
			name:     "Error/BadTarget",
			tokenRaw: "token",
			form: models.VoteForm{
				TargetID: test.NumberUUID(1),
				Target:   "fake-target",
				Vote:     lo.ToPtr(models.VoteValueUp),
			},
			id:         test.NumberUUID(10),
			now:        baseTime,
			clientName: "target",
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
			form: models.VoteForm{
				TargetID: test.NumberUUID(1),
				Target:   "fake-target",
				Vote:     lo.ToPtr(models.VoteValueUp),
			},
			id:             test.NumberUUID(10),
			now:            baseTime,
			clientName:     "target",
			authClientResp: &authmodels.UserTokenStatus{},
			expectErr:      errors.ErrInvalidCredentials,
		},
		{
			name:     "Error/AuthClientFailure",
			tokenRaw: "token",
			form: models.VoteForm{
				TargetID: test.NumberUUID(1),
				Target:   "fake-target",
				Vote:     lo.ToPtr(models.VoteValueUp),
			},
			id:            test.NumberUUID(10),
			now:           baseTime,
			clientName:    "target",
			authClientErr: fooErr,
			expectErr:     fooErr,
		},
		{
			name:     "Error/BadVote",
			tokenRaw: "token",
			form: models.VoteForm{
				TargetID: test.NumberUUID(1),
				Target:   "target",
				Vote:     lo.ToPtr(models.VoteValue("invalid")),
			},
			id:         test.NumberUUID(10),
			now:        baseTime,
			clientName: "target",
			authClientResp: &authmodels.UserTokenStatus{
				OK: true,
				Token: &authmodels.UserToken{
					Payload: authmodels.UserTokenPayload{ID: test.NumberUUID(100)},
				},
			},
			expectErr: errors.ErrInvalidEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			repository := daomocks.NewVotesRepository(t)
			authClient := authmocks.NewClient(t)

			authClient.On("IntrospectToken", context.Background(), d.tokenRaw).Return(d.authClientResp, d.authClientErr)

			if d.shouldCallDAO {
				repository.
					On("Cast", context.Background(), d.authClientResp.Token.Payload.ID, d.form.TargetID, d.form.Target, d.form.Vote, d.id, d.now).
					Return(nil, d.castErr)

				// Execute the actual method, but call the mocks inside of it.
				txCall := repository.On("RunInTx", context.Background(), mock.Anything)
				txCall.Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context, dao.VotesRepository) error)
					txCall.ReturnArguments = []interface{}{fn(context.Background(), repository)}
				})
			}

			if d.shouldCallGetSummary {
				repository.
					On("GetSummary", context.Background(), d.form.TargetID, d.form.Target).
					Return(d.summary, d.summaryErr)
			}

			targets := map[string]models.CheckVoteClient{
				d.clientName: func(ctx context.Context, id, userID uuid.UUID, upVotes, downVotes int) error {
					return d.clientErr
				},
			}

			service := services.NewCastVoteService(repository, authClient, targets)
			res, err := service.Cast(context.Background(), d.tokenRaw, d.form, d.id, d.now)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, res)

			repository.AssertExpectations(t)
			authClient.AssertExpectations(t)
		})
	}
}
