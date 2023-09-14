package dao

import (
	"context"
	"github.com/a-novel/bunovel"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type VotesRepository interface {
	Get(ctx context.Context, userID, targetID uuid.UUID, target string) (*VoteModel, error)
	GetSummary(ctx context.Context, targetID uuid.UUID, target string) (*VotesSummaryModel, error)
	ListUserVotes(ctx context.Context, userID uuid.UUID, target string, limit, offset int) ([]*VoteModel, error)
	Cast(ctx context.Context, userID, targetID uuid.UUID, target string, vote *models.VoteValue, id uuid.UUID, now time.Time) (*VoteModel, error)

	RunInTx(ctx context.Context, f func(ctx context.Context, txClient VotesRepository) error) error
}

type VoteModel struct {
	bun.BaseModel `bun:"table:votes"`
	bunovel.Metadata

	Vote     models.VoteValue `bun:"type:vote"`
	UserID   uuid.UUID        `bun:"user_id"`
	TargetID uuid.UUID        `bun:"target_id"`
	Target   string           `bun:"target"`
}

type VotesSummaryModel struct {
	bun.BaseModel `bun:"table:votes_summary"`

	TargetID  uuid.UUID `bun:"target_id"`
	Target    string    `bun:"target"`
	UpVotes   int       `bun:"up_votes"`
	DownVotes int       `bun:"down_votes"`
}

func NewVotesRepository(db bun.IDB) VotesRepository {
	return &votesRepositoryImpl{db: db}
}

type votesRepositoryImpl struct {
	db bun.IDB
}

func (repository *votesRepositoryImpl) Get(ctx context.Context, userID, targetID uuid.UUID, target string) (*VoteModel, error) {
	model := new(VoteModel)

	err := repository.db.NewSelect().Model(model).
		Where("user_id = ?", userID).
		Where("target_id = ?", targetID).
		Where("target = ?", target).
		Scan(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	return model, nil
}

func (repository *votesRepositoryImpl) GetSummary(ctx context.Context, targetID uuid.UUID, target string) (*VotesSummaryModel, error) {
	model := new(VotesSummaryModel)

	err := repository.db.NewSelect().Model(model).
		Where("target_id = ?", targetID).
		Where("target = ?", target).
		Scan(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	return model, nil
}

func (repository *votesRepositoryImpl) ListUserVotes(ctx context.Context, userID uuid.UUID, target string, limit, offset int) ([]*VoteModel, error) {
	votes := make([]*VoteModel, 0)

	err := repository.db.NewSelect().Model(&votes).
		Where("user_id = ?", userID).
		Where("target = ?", target).
		OrderExpr("COALESCE(updated_at, created_at) DESC").
		Limit(limit).Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	return votes, nil
}

func (repository *votesRepositoryImpl) Cast(ctx context.Context, userID, targetID uuid.UUID, target string, vote *models.VoteValue, id uuid.UUID, now time.Time) (*VoteModel, error) {
	model := new(VoteModel)

	if vote == nil {
		_, err := repository.db.NewDelete().Model(model).
			Where("user_id = ?", userID).
			Where("target_id = ?", targetID).
			Where("target = ?", target).
			Exec(ctx)

		if err != nil {
			return nil, bunovel.HandlePGError(err)
		}

		return nil, nil
	}

	model.ID = id
	model.CreatedAt = now
	model.UserID = userID
	model.TargetID = targetID
	model.Target = target
	model.Vote = *vote

	err := repository.db.NewInsert().Model(model).
		Returning("*").
		On("CONFLICT (user_id, target_id, target) DO UPDATE").
		Set("vote = EXCLUDED.vote").
		Set("updated_at = EXCLUDED.created_at").
		Scan(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	return model, nil
}

func (repository *votesRepositoryImpl) RunInTx(ctx context.Context, callback func(ctx context.Context, txRepository VotesRepository) error) error {
	return repository.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return callback(ctx, NewVotesRepository(tx))
	})
}
