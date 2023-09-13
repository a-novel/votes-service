package services

import (
	goerrors "errors"
)

var (
	ErrInvalidToken       = goerrors.New("(data) invalid tokenRaw")
	ErrInvalidSearchLimit = goerrors.New("(data) invalid search limit")
	ErrInvalidTarget      = goerrors.New("(data) invalid target")

	ErrIntrospectToken  = goerrors.New("(dep) failed to introspect tokenRaw")
	ErrSendVoteToTarget = goerrors.New("(dep) failed to send vote to target")

	ErrGetVote         = goerrors.New("(dao) failed to get vote")
	ErrListUserVotes   = goerrors.New("(dao) failed to list user votes")
	ErrCastVote        = goerrors.New("(dao) failed to cast vote")
	ErrGetVotesSummary = goerrors.New("(dao) failed to get votes summary")
)

const (
	MaxSearchLimit = 100
)
