package dao_test

import (
	"context"
	"github.com/a-novel/bunovel"
	goframework "github.com/a-novel/go-framework"
	"github.com/a-novel/votes-service/migrations"
	"github.com/a-novel/votes-service/pkg/dao"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"io/fs"
	"testing"
	"time"
)

func TestVotesRepository_Get(t *testing.T) {
	db, sqlDB := bunovel.GetTestPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.VoteModel{
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   goframework.NumberUUID(1),
			TargetID: goframework.NumberUUID(1),
			Target:   "target",
		},
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(3), baseTime, nil),
			Vote:     models.VoteValueDown,
			UserID:   goframework.NumberUUID(3),
			TargetID: goframework.NumberUUID(1),
			Target:   "target",
		},
	}

	data := []struct {
		name string

		userID   uuid.UUID
		targetID uuid.UUID
		target   string

		expect    *dao.VoteModel
		expectErr error
	}{
		{
			name:     "Success",
			userID:   goframework.NumberUUID(3),
			targetID: goframework.NumberUUID(1),
			target:   "target",
			expect: &dao.VoteModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(3), baseTime, nil),
				Vote:     models.VoteValueDown,
				UserID:   goframework.NumberUUID(3),
				TargetID: goframework.NumberUUID(1),
				Target:   "target",
			},
		},
		{
			name:      "Error/NotFound",
			userID:    goframework.NumberUUID(3),
			targetID:  goframework.NumberUUID(1),
			target:    "other-target",
			expectErr: bunovel.ErrNotFound,
		},
	}

	err := bunovel.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := dao.NewVotesRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.Get(ctx, d.userID, d.targetID, d.target)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestVotesRepository_GetSummary(t *testing.T) {
	db, sqlDB := bunovel.GetTestPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.VoteModel{
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   goframework.NumberUUID(1),
			TargetID: goframework.NumberUUID(1),
			Target:   "target",
		},
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(2), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   goframework.NumberUUID(2),
			TargetID: goframework.NumberUUID(1),
			Target:   "target",
		},
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(3), baseTime, nil),
			Vote:     models.VoteValueDown,
			UserID:   goframework.NumberUUID(3),
			TargetID: goframework.NumberUUID(1),
			Target:   "target",
		},

		// Another target id.
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(4), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   goframework.NumberUUID(1),
			TargetID: goframework.NumberUUID(2),
			Target:   "target",
		},

		// Another target.
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(5), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   goframework.NumberUUID(1),
			TargetID: goframework.NumberUUID(1),
			Target:   "other-target",
		},
	}

	data := []struct {
		name string

		targetID uuid.UUID
		target   string

		expect    *dao.VotesSummaryModel
		expectErr error
	}{
		{
			name:     "Success",
			targetID: goframework.NumberUUID(1),
			target:   "target",
			expect: &dao.VotesSummaryModel{
				Target:    "target",
				TargetID:  goframework.NumberUUID(1),
				UpVotes:   2,
				DownVotes: 1,
			},
		},
		{
			name:      "Error/NotFound",
			targetID:  goframework.NumberUUID(10),
			target:    "target",
			expectErr: bunovel.ErrNotFound,
		},
	}

	err := bunovel.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := dao.NewVotesRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.GetSummary(ctx, d.targetID, d.target)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestVotesRepository_ListUserVotes(t *testing.T) {
	db, sqlDB := bunovel.GetTestPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.VoteModel{
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   goframework.NumberUUID(1),
			TargetID: goframework.NumberUUID(1),
			Target:   "target",
		},
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(2), baseTime.Add(30*time.Minute), nil),
			Vote:     models.VoteValueUp,
			UserID:   goframework.NumberUUID(2),
			TargetID: goframework.NumberUUID(1),
			Target:   "target",
		},

		// Another target id.
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(4), baseTime, lo.ToPtr(updateTime.Add(time.Hour))),
			Vote:     models.VoteValueUp,
			UserID:   goframework.NumberUUID(2),
			TargetID: goframework.NumberUUID(2),
			Target:   "target",
		},

		// Another target.
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(5), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   goframework.NumberUUID(2),
			TargetID: goframework.NumberUUID(1),
			Target:   "other-target",
		},
	}

	data := []struct {
		name string

		userID uuid.UUID
		target string
		limit  int
		offset int

		expect    []*dao.VoteModel
		expectErr error
	}{
		{
			name:   "Success",
			userID: goframework.NumberUUID(2),
			target: "target",
			expect: []*dao.VoteModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(4), baseTime, lo.ToPtr(updateTime.Add(time.Hour))),
					Vote:     models.VoteValueUp,
					UserID:   goframework.NumberUUID(2),
					TargetID: goframework.NumberUUID(2),
					Target:   "target",
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(2), baseTime.Add(30*time.Minute), nil),
					Vote:     models.VoteValueUp,
					UserID:   goframework.NumberUUID(2),
					TargetID: goframework.NumberUUID(1),
					Target:   "target",
				},
			},
		},
		{
			name:   "Success/Limit",
			userID: goframework.NumberUUID(2),
			target: "target",
			limit:  1,
			expect: []*dao.VoteModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(4), baseTime, lo.ToPtr(updateTime.Add(time.Hour))),
					Vote:     models.VoteValueUp,
					UserID:   goframework.NumberUUID(2),
					TargetID: goframework.NumberUUID(2),
					Target:   "target",
				},
			},
		},
		{
			name:   "Success/Offset",
			userID: goframework.NumberUUID(2),
			target: "target",
			offset: 1,
			expect: []*dao.VoteModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(2), baseTime.Add(30*time.Minute), nil),
					Vote:     models.VoteValueUp,
					UserID:   goframework.NumberUUID(2),
					TargetID: goframework.NumberUUID(1),
					Target:   "target",
				},
			},
		},
		{
			name:   "Success/NoResults",
			userID: goframework.NumberUUID(10),
			target: "target",
			expect: []*dao.VoteModel{},
		},
	}

	err := bunovel.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := dao.NewVotesRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.ListUserVotes(ctx, d.userID, d.target, d.limit, d.offset)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestVotesRepository_Cast(t *testing.T) {
	db, sqlDB := bunovel.GetTestPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.VoteModel{
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   goframework.NumberUUID(1),
			TargetID: goframework.NumberUUID(1),
			Target:   "target",
		},
	}

	data := []struct {
		name string

		userID   uuid.UUID
		targetID uuid.UUID
		target   string
		vote     *models.VoteValue
		id       uuid.UUID
		now      time.Time

		expect    *dao.VoteModel
		expectErr error
	}{
		{
			name:     "Success",
			userID:   goframework.NumberUUID(2),
			targetID: goframework.NumberUUID(1),
			target:   "target",
			vote:     lo.ToPtr(models.VoteValueDown),
			id:       goframework.NumberUUID(2),
			now:      updateTime,
			expect: &dao.VoteModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(2), updateTime, nil),
				Vote:     models.VoteValueDown,
				UserID:   goframework.NumberUUID(2),
				TargetID: goframework.NumberUUID(1),
				Target:   "target",
			},
		},
		{
			name:     "Success/Update",
			userID:   goframework.NumberUUID(1),
			targetID: goframework.NumberUUID(1),
			target:   "target",
			vote:     lo.ToPtr(models.VoteValueDown),
			id:       goframework.NumberUUID(2),
			now:      updateTime,
			expect: &dao.VoteModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &updateTime),
				Vote:     models.VoteValueDown,
				UserID:   goframework.NumberUUID(1),
				TargetID: goframework.NumberUUID(1),
				Target:   "target",
			},
		},
		{
			name:     "Success/Delete",
			userID:   goframework.NumberUUID(1),
			targetID: goframework.NumberUUID(1),
			target:   "target",
			id:       goframework.NumberUUID(2),
			now:      updateTime,
		},
		{
			name:     "Success/DeleteMissing",
			userID:   goframework.NumberUUID(2),
			targetID: goframework.NumberUUID(1),
			target:   "target",
			id:       goframework.NumberUUID(2),
			now:      updateTime,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			err := bunovel.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
				repository := dao.NewVotesRepository(tx)

				res, err := repository.Cast(ctx, d.userID, d.targetID, d.target, d.vote, d.id, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
			require.NoError(t, err)
		})
	}
}
