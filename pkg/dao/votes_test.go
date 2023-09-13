package dao_test

import (
	"context"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/postgresql"
	"github.com/a-novel/go-framework/test"
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
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.VoteModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   test.NumberUUID(1),
			TargetID: test.NumberUUID(1),
			Target:   "target",
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(3), baseTime, nil),
			Vote:     models.VoteValueDown,
			UserID:   test.NumberUUID(3),
			TargetID: test.NumberUUID(1),
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
			userID:   test.NumberUUID(3),
			targetID: test.NumberUUID(1),
			target:   "target",
			expect: &dao.VoteModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(3), baseTime, nil),
				Vote:     models.VoteValueDown,
				UserID:   test.NumberUUID(3),
				TargetID: test.NumberUUID(1),
				Target:   "target",
			},
		},
		{
			name:      "Error/NotFound",
			userID:    test.NumberUUID(3),
			targetID:  test.NumberUUID(1),
			target:    "other-target",
			expectErr: errors.ErrNotFound,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
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
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.VoteModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   test.NumberUUID(1),
			TargetID: test.NumberUUID(1),
			Target:   "target",
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(2), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   test.NumberUUID(2),
			TargetID: test.NumberUUID(1),
			Target:   "target",
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(3), baseTime, nil),
			Vote:     models.VoteValueDown,
			UserID:   test.NumberUUID(3),
			TargetID: test.NumberUUID(1),
			Target:   "target",
		},

		// Another target id.
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(4), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   test.NumberUUID(1),
			TargetID: test.NumberUUID(2),
			Target:   "target",
		},

		// Another target.
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(5), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   test.NumberUUID(1),
			TargetID: test.NumberUUID(1),
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
			targetID: test.NumberUUID(1),
			target:   "target",
			expect: &dao.VotesSummaryModel{
				Target:    "target",
				TargetID:  test.NumberUUID(1),
				UpVotes:   2,
				DownVotes: 1,
			},
		},
		{
			name:      "Error/NotFound",
			targetID:  test.NumberUUID(10),
			target:    "target",
			expectErr: errors.ErrNotFound,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
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
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.VoteModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   test.NumberUUID(1),
			TargetID: test.NumberUUID(1),
			Target:   "target",
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(2), baseTime.Add(30*time.Minute), nil),
			Vote:     models.VoteValueUp,
			UserID:   test.NumberUUID(2),
			TargetID: test.NumberUUID(1),
			Target:   "target",
		},

		// Another target id.
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(4), baseTime, lo.ToPtr(updateTime.Add(time.Hour))),
			Vote:     models.VoteValueUp,
			UserID:   test.NumberUUID(2),
			TargetID: test.NumberUUID(2),
			Target:   "target",
		},

		// Another target.
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(5), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   test.NumberUUID(2),
			TargetID: test.NumberUUID(1),
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
			userID: test.NumberUUID(2),
			target: "target",
			expect: []*dao.VoteModel{
				{
					Metadata: postgresql.NewMetadata(test.NumberUUID(4), baseTime, lo.ToPtr(updateTime.Add(time.Hour))),
					Vote:     models.VoteValueUp,
					UserID:   test.NumberUUID(2),
					TargetID: test.NumberUUID(2),
					Target:   "target",
				},
				{
					Metadata: postgresql.NewMetadata(test.NumberUUID(2), baseTime.Add(30*time.Minute), nil),
					Vote:     models.VoteValueUp,
					UserID:   test.NumberUUID(2),
					TargetID: test.NumberUUID(1),
					Target:   "target",
				},
			},
		},
		{
			name:   "Success/Limit",
			userID: test.NumberUUID(2),
			target: "target",
			limit:  1,
			expect: []*dao.VoteModel{
				{
					Metadata: postgresql.NewMetadata(test.NumberUUID(4), baseTime, lo.ToPtr(updateTime.Add(time.Hour))),
					Vote:     models.VoteValueUp,
					UserID:   test.NumberUUID(2),
					TargetID: test.NumberUUID(2),
					Target:   "target",
				},
			},
		},
		{
			name:   "Success/Offset",
			userID: test.NumberUUID(2),
			target: "target",
			offset: 1,
			expect: []*dao.VoteModel{
				{
					Metadata: postgresql.NewMetadata(test.NumberUUID(2), baseTime.Add(30*time.Minute), nil),
					Vote:     models.VoteValueUp,
					UserID:   test.NumberUUID(2),
					TargetID: test.NumberUUID(1),
					Target:   "target",
				},
			},
		},
		{
			name:   "Success/NoResults",
			userID: test.NumberUUID(10),
			target: "target",
			expect: []*dao.VoteModel{},
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
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
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.VoteModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, nil),
			Vote:     models.VoteValueUp,
			UserID:   test.NumberUUID(1),
			TargetID: test.NumberUUID(1),
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
			userID:   test.NumberUUID(2),
			targetID: test.NumberUUID(1),
			target:   "target",
			vote:     lo.ToPtr(models.VoteValueDown),
			id:       test.NumberUUID(2),
			now:      updateTime,
			expect: &dao.VoteModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(2), updateTime, nil),
				Vote:     models.VoteValueDown,
				UserID:   test.NumberUUID(2),
				TargetID: test.NumberUUID(1),
				Target:   "target",
			},
		},
		{
			name:     "Success/Update",
			userID:   test.NumberUUID(1),
			targetID: test.NumberUUID(1),
			target:   "target",
			vote:     lo.ToPtr(models.VoteValueDown),
			id:       test.NumberUUID(2),
			now:      updateTime,
			expect: &dao.VoteModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, &updateTime),
				Vote:     models.VoteValueDown,
				UserID:   test.NumberUUID(1),
				TargetID: test.NumberUUID(1),
				Target:   "target",
			},
		},
		{
			name:     "Success/Delete",
			userID:   test.NumberUUID(1),
			targetID: test.NumberUUID(1),
			target:   "target",
			id:       test.NumberUUID(2),
			now:      updateTime,
		},
		{
			name:     "Success/DeleteMissing",
			userID:   test.NumberUUID(2),
			targetID: test.NumberUUID(1),
			target:   "target",
			id:       test.NumberUUID(2),
			now:      updateTime,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
				repository := dao.NewVotesRepository(tx)

				res, err := repository.Cast(ctx, d.userID, d.targetID, d.target, d.vote, d.id, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
			require.NoError(t, err)
		})
	}
}
