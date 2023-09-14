package services

import (
	"context"
	goerrors "errors"
	auth "github.com/a-novel/auth-service/framework"
	goframework "github.com/a-novel/go-framework"
	"github.com/a-novel/votes-service/pkg/adapters"
	"github.com/a-novel/votes-service/pkg/dao"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"time"
)

type CastVoteService interface {
	Cast(ctx context.Context, tokenRaw string, form models.VoteForm, id uuid.UUID, now time.Time) (*models.VotesSummary, error)
}

func NewCastVoteService(repository dao.VotesRepository, authClient auth.Client, targetsClients map[string]models.CheckVoteClient) CastVoteService {
	return &castVoteServiceImpl{
		repository:     repository,
		authClient:     authClient,
		targetsClients: targetsClients,
	}
}

type castVoteServiceImpl struct {
	repository dao.VotesRepository
	authClient auth.Client

	targetsClients map[string]models.CheckVoteClient
}

func (s *castVoteServiceImpl) Cast(ctx context.Context, tokenRaw string, form models.VoteForm, id uuid.UUID, now time.Time) (*models.VotesSummary, error) {
	var res *dao.VotesSummaryModel

	token, err := s.authClient.IntrospectToken(ctx, tokenRaw)
	if err != nil {
		return nil, goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return nil, goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidToken)
	}

	if err := goframework.CheckRestricted(lo.FromPtr(form.Vote), "", models.VoteValueUp, models.VoteValueDown); err != nil {
		return nil, goerrors.Join(goframework.ErrInvalidEntity, err)
	}

	targetClient := s.targetsClients[form.Target]
	if targetClient == nil {
		return nil, goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidTarget)
	}

	// Prevent insertion if client call fails.
	err = s.repository.RunInTx(ctx, func(ctx context.Context, txRepository dao.VotesRepository) error {
		_, err = txRepository.Cast(ctx, token.Token.Payload.ID, form.TargetID, form.Target, form.Vote, id, now)
		if err != nil {
			return goerrors.Join(ErrCastVote, err)
		}

		res, err = txRepository.GetSummary(ctx, form.TargetID, form.Target)
		if err != nil {
			return goerrors.Join(ErrGetVotesSummary, err)
		}

		if err := targetClient(ctx, form.TargetID, token.Token.Payload.ID, res.UpVotes, res.DownVotes); err != nil {
			return goerrors.Join(ErrSendVoteToTarget, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return adapters.VotesSummaryToModel(res), nil
}
