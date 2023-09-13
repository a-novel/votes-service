package services_test

import (
	"context"
	"github.com/a-novel/go-framework/test"
	"github.com/a-novel/votes-service/pkg/dao"
	daomocks "github.com/a-novel/votes-service/pkg/dao/mocks"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/a-novel/votes-service/pkg/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetVotesSummaryService(t *testing.T) {
	data := []struct {
		name string

		targetID uuid.UUID
		target   string

		daoResp *dao.VotesSummaryModel
		daoErr  error

		expect    *models.VotesSummary
		expectErr error
	}{
		{
			name:     "Success",
			targetID: test.NumberUUID(1),
			target:   "target",
			daoResp: &dao.VotesSummaryModel{
				UpVotes:   100,
				DownVotes: 50,
			},
			expect: &models.VotesSummary{
				UpVotes:   100,
				DownVotes: 50,
			},
		},
		{
			name:      "Error/DAOFailure",
			targetID:  test.NumberUUID(1),
			target:    "target",
			daoErr:    fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			repository := daomocks.NewVotesRepository(t)

			repository.On("GetSummary", context.Background(), d.targetID, d.target).Return(d.daoResp, d.daoErr)

			service := services.NewGetVotesSummaryService(repository)
			resp, err := service.Get(context.Background(), d.targetID, d.target)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, resp)

			repository.AssertExpectations(t)
		})
	}
}
