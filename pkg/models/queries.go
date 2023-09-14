package models

import (
	"github.com/a-novel/go-apis"
)

type ListUserVotesQuery struct {
	Target string `json:"target" form:"target"`
	Limit  int    `json:"limit" form:"limit"`
	Offset int    `json:"offset" form:"offset"`
}

type GetUserVoteQuery struct {
	TargetID apis.StringUUID `json:"targetID" form:"targetID"`
	Target   string          `json:"target" form:"target"`
}

type GetVotesSummaryQuery struct {
	TargetID apis.StringUUID `json:"targetID" form:"targetID"`
	Target   string          `json:"target" form:"target"`
}
