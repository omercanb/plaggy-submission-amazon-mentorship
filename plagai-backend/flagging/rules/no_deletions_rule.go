package rules

import (
	"strings"

	"github.com/plagai/plagai-backend/models/domain"
)

type NoDeletionsRule struct{}

func (r NoDeletionsRule) Apply(diffs []domain.Diff) []domain.Flag {
	for _, diff := range diffs {
		if strings.HasPrefix(diff.PatchText, "-") {
			return nil // deletion found, no flag
		}
	}

	return []domain.Flag{
		{
			Diff: domain.Diff{},

			FlagExplanation: "No deletions in the assignment",
			Severity:        1,
		},
	}
}
