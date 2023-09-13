package models

import "github.com/google/uuid"

type VoteForm struct {
	TargetID uuid.UUID  `json:"targetID" form:"targetID"`
	Target   string     `json:"target" form:"target"`
	Vote     *VoteValue `json:"vote" form:"vote"`
}
