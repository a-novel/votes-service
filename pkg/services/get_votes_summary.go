package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/votes-service/pkg/adapters"
	"github.com/a-novel/votes-service/pkg/dao"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/google/uuid"
)

type GetVotesSummaryService interface {
	Get(ctx context.Context, targetID uuid.UUID, target string) (*models.VotesSummary, error)
}

func NewGetVotesSummaryService(repository dao.VotesRepository) GetVotesSummaryService {
	return &getVotesSummaryServiceImpl{
		repository: repository,
	}
}

type getVotesSummaryServiceImpl struct {
	repository dao.VotesRepository
}

func (s *getVotesSummaryServiceImpl) Get(ctx context.Context, targetID uuid.UUID, target string) (*models.VotesSummary, error) {
	summary, err := s.repository.GetSummary(ctx, targetID, target)
	if err != nil {
		return nil, goerrors.Join(ErrGetVotesSummary, err)
	}

	return adapters.VotesSummaryToModel(summary), nil
}
