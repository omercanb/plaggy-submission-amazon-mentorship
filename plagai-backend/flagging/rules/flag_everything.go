package rules

import (
	"time"

	"github.com/plagai/plagai-backend/models/domain"
)

type FlagEverythingRule struct{}

func (r FlagEverythingRule) Apply(diff domain.Diff, prevTimestamp time.Time) *domain.Flag {
	return &domain.Flag{
		Diff:            diff,
		FlagExplanation: "Flagging everything for testing",
		Severity:        1,
	}
}
