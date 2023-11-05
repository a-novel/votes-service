package services

import (
	"context"
	goerrors "errors"
	apiclients "github.com/a-novel/go-apis/clients"
	goframework "github.com/a-novel/go-framework"
	"github.com/a-novel/votes-service/pkg/adapters"
	"github.com/a-novel/votes-service/pkg/dao"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/google/uuid"
)

type GetUserVoteService interface {
	Get(ctx context.Context, tokenRaw string, targetID uuid.UUID, target string) (*models.Vote, error)
}

func NewGetUserVoteService(repository dao.VotesRepository, authClient apiclients.AuthClient) GetUserVoteService {
	return &getUserVoteServiceImpl{
		repository: repository,
		authClient: authClient,
	}
}

type getUserVoteServiceImpl struct {
	repository dao.VotesRepository
	authClient apiclients.AuthClient
}

func (s *getUserVoteServiceImpl) Get(ctx context.Context, tokenRaw string, targetID uuid.UUID, target string) (*models.Vote, error) {
	token, err := s.authClient.IntrospectToken(ctx, tokenRaw)
	if err != nil {
		return nil, goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return nil, goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidToken)
	}

	vote, err := s.repository.Get(ctx, token.Token.Payload.ID, targetID, target)
	if err != nil {
		return nil, goerrors.Join(ErrGetVote, err)
	}

	return adapters.VoteToModel(vote), nil
}
