package flagging

import (
	"time"

	"github.com/plagai/plagai-backend/flagging/rules"
	"github.com/plagai/plagai-backend/models"
	"github.com/plagai/plagai-backend/models/domain"
)

type FlaggingEngine struct {
	DiffRules       []DiffRule
	AssignmentRules []AssignmentRule
}

// Rule applied to a single diff
type DiffRule interface {
	Apply(diff domain.Diff, prevTimestamp time.Time) *domain.Flag
}

// Rule applied to the full assignment (all diffs)
type AssignmentRule interface {
	Apply(diffs []domain.Diff) []domain.Flag
}

func GetDefaultFlaggingEngine() *FlaggingEngine {
	return &FlaggingEngine{
		DiffRules: []DiffRule{
			rules.SpeedThresholdRule{MaxCharsPerSecond: 20}, // example threshold
			rules.FlagEverythingRule{},
		},
		AssignmentRules: []AssignmentRule{
			rules.NoDeletionsRule{},
		},
	}
}

func NewFlaggingEngine(diffRules []DiffRule, assignmentRules []AssignmentRule) *FlaggingEngine {
	return &FlaggingEngine{
		DiffRules:       diffRules,
		AssignmentRules: assignmentRules,
	}
}

func (e *FlaggingEngine) FlagAssignment(events []models.EditEvent) []domain.Flag {
	flags := []domain.Flag{}

	// Accumulate diffs in the diff rules to use for assignment rules
	lastEditTimeForFile := make(map[string]time.Time)
	diffs := []domain.Diff{}
	// Apply per-diff rules
	for _, event := range events {
		if event.EventType != models.APIEventAdded && event.EventType != models.APIEventModified {
			lastEditTimeForFile[event.FilePath] = event.Timestamp
			continue
		}
		diff := domain.Diff{
			FilePath:  event.FilePath,
			PatchText: event.Patch,
			Timestamp: event.Timestamp,
		}
		diffs = append(diffs, diff)
		for _, rule := range e.DiffRules {
			prevEditTime := lastEditTimeForFile[event.FilePath]
			if f := rule.Apply(diff, prevEditTime); f != nil {
				f.Diff = diff
				flags = append(flags, *f)
			}
		}
		lastEditTimeForFile[event.FilePath] = event.Timestamp
	}

	// Apply whole-assignment rules
	for _, rule := range e.AssignmentRules {
		flags = append(flags, rule.Apply(diffs)...)
	}

	return flags
}
