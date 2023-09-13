package models

import (
	"github.com/google/uuid"
	"time"
)

type VoteValue string

var (
	VoteValueUp   VoteValue = "up"
	VoteValueDown VoteValue = "down"
)

type Vote struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updatedAt"`

	Vote     VoteValue `json:"vote"`
	UserID   uuid.UUID `json:"userID"`
	TargetID uuid.UUID `json:"targetID"`
	Target   string    `json:"target"`
}

type VotesSummary struct {
	UpVotes   int `json:"upVotes"`
	DownVotes int `json:"downVotes"`
}
