package adapters

import (
	"github.com/a-novel/votes-service/pkg/dao"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/samber/lo"
)

func VoteToModel(src *dao.VoteModel) *models.Vote {
	if src == nil {
		return nil
	}

	return &models.Vote{
		ID:        src.ID,
		Vote:      src.Vote,
		UserID:    src.UserID,
		TargetID:  src.TargetID,
		Target:    src.Target,
		UpdatedAt: lo.Ternary(src.UpdatedAt == nil, src.CreatedAt, lo.FromPtr(src.UpdatedAt)),
	}
}
