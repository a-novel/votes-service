package adapters

import (
	"context"
	apiclients "github.com/a-novel/go-apis/clients"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/google/uuid"
)

func NewImproveRequestVoteClient(client apiclients.ForumClient, permissionsClient apiclients.PermissionsClient) models.CheckVoteClient {
	return func(ctx context.Context, id, userID uuid.UUID, upVotes, downVotes int) error {
		if err := permissionsClient.HasUserScope(ctx, apiclients.HasUserScopeQuery{
			UserID: userID,
			Scope:  apiclients.CanVotePost,
		}); err != nil {
			return err
		}

		return client.VoteImproveRequest(ctx, apiclients.UpdateImproveRequestVotesForm{
			ID:        id,
			UserID:    userID,
			UpVotes:   upVotes,
			DownVotes: downVotes,
		})
	}
}

func NewImproveSuggestionVoteClient(client apiclients.ForumClient, permissionsClient apiclients.PermissionsClient) models.CheckVoteClient {
	return func(ctx context.Context, id, userID uuid.UUID, upVotes, downVotes int) error {
		if err := permissionsClient.HasUserScope(ctx, apiclients.HasUserScopeQuery{
			UserID: userID,
			Scope:  apiclients.CanVotePost,
		}); err != nil {
			return err
		}

		return client.VoteImproveSuggestion(ctx, apiclients.UpdateImproveSuggestionVotesForm{
			ID:        id,
			UserID:    userID,
			UpVotes:   upVotes,
			DownVotes: downVotes,
		})
	}
}
