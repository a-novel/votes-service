package services

import (
	"context"
	goerrors "errors"
	apiclients "github.com/a-novel/go-apis/clients"
	goframework "github.com/a-novel/go-framework"
	"github.com/a-novel/votes-service/pkg/adapters"
	"github.com/a-novel/votes-service/pkg/dao"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/samber/lo"
)

type ListUserVotesService interface {
	List(ctx context.Context, tokenRaw string, query *models.ListUserVotesQuery) ([]*models.Vote, error)
}

func NewListUserVotesService(repository dao.VotesRepository, authClient apiclients.AuthClient) ListUserVotesService {
	return &listUserVotesServiceImpl{
		repository: repository,
		authClient: authClient,
	}
}

type listUserVotesServiceImpl struct {
	repository dao.VotesRepository
	authClient apiclients.AuthClient
}

func (s *listUserVotesServiceImpl) List(ctx context.Context, tokenRaw string, query *models.ListUserVotesQuery) ([]*models.Vote, error) {
	token, err := s.authClient.IntrospectToken(ctx, tokenRaw)
	if err != nil {
		return nil, goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return nil, goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidToken)
	}

	if err := goframework.CheckMinMax(query.Limit, 1, MaxSearchLimit); err != nil {
		return nil, goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidSearchLimit, err)
	}

	votes, err := s.repository.ListUserVotes(ctx, token.Token.Payload.ID, query.Target, query.Limit, query.Offset)
	if err != nil {
		return nil, goerrors.Join(ErrListUserVotes, err)
	}

	return lo.Map(votes, func(item *dao.VoteModel, _ int) *models.Vote {
		return adapters.VoteToModel(item)
	}), nil
}
