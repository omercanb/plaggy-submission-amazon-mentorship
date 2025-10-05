package rules

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/plagai/plagai-backend/models/domain"
)

type SpeedThresholdRule struct {
	MaxCharsPerSecond float64 // e.g., 5–10 chars/sec is reasonable, >50 chars/sec is suspicious
}

func (r SpeedThresholdRule) Apply(diff domain.Diff, prevTimestamp time.Time) *domain.Flag {
	if prevTimestamp.IsZero() {
		// First diff, cannot compute speed, just return nil
		return nil
	}

	duration := diff.Timestamp.Sub(prevTimestamp).Seconds()
	if duration <= 0 {
		// Duration should never be less than zero
		log.Printf("problematic duration calculation <= 0, for diff with filepath: %v, timestamp: %v\n", diff.FilePath, diff.Timestamp)
		return nil
	}

	lengthOfAdditions := 0
	for line := range strings.SplitSeq(diff.PatchText, "\n") {
		if len(line) != 0 && line[0] == '+' {
			lengthOfAdditions += len(line) - 1
		}
	}
	speed := float64(lengthOfAdditions) / duration

	if speed > r.MaxCharsPerSecond {
		return &domain.Flag{
			Diff: diff,
			FlagExplanation: "Text entered too fast (probably copy-pasted), speed: " +
				fmt.Sprintf("%.1f chars/sec", speed),
			Severity: 2,
		}
	}

	return nil
}
