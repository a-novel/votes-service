package models

import (
	"context"
	"github.com/google/uuid"
)

type CheckVoteClient func(ctx context.Context, id, userID uuid.UUID, upVotes, downVotes int) error
