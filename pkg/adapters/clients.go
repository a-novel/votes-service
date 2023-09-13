package adapters

import (
	"context"
	forumclient "github.com/a-novel/forum-service/framework"
	forummodels "github.com/a-novel/forum-service/pkg/models"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/google/uuid"
)

func NewImproveRequestVoteClient(client forumclient.Client) models.CheckVoteClient {
	return func(ctx context.Context, id, userID uuid.UUID, upVotes, downVotes int) error {
		return client.VoteImproveRequest(ctx, forummodels.UpdateImproveRequestVotesForm{
			ID:        id,
			UserID:    userID,
			UpVotes:   upVotes,
			DownVotes: downVotes,
		})
	}
}

func NewImproveSuggestionVoteClient(client forumclient.Client) models.CheckVoteClient {
	return func(ctx context.Context, id, userID uuid.UUID, upVotes, downVotes int) error {
		return client.VoteImproveSuggestion(ctx, forummodels.UpdateImproveSuggestionVotesForm{
			ID:        id,
			UserID:    userID,
			UpVotes:   upVotes,
			DownVotes: downVotes,
		})
	}
}
