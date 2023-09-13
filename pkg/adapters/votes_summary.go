package adapters

import (
	"github.com/a-novel/votes-service/pkg/dao"
	"github.com/a-novel/votes-service/pkg/models"
)

func VotesSummaryToModel(src *dao.VotesSummaryModel) *models.VotesSummary {
	if src == nil {
		return nil
	}

	return &models.VotesSummary{
		UpVotes:   src.UpVotes,
		DownVotes: src.DownVotes,
	}
}
