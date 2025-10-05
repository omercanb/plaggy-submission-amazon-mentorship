package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/plagai/plagai-backend/models/domain"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func BuildFilesystemFromPatches(patches []domain.Diff) (map[string]string, error) {
	filepathToPatches := make(map[string][]domain.Diff)
	for _, patch := range patches {
		filepathToPatches[patch.FilePath] = append(filepathToPatches[patch.FilePath], patch)
	}
	result := make(map[string]string)
	for filepath, patchesForFile := range filepathToPatches {
		fileFinalState, err := BuildFileFromPatches(patchesForFile)
		if err != nil {
			log.Println(err)
		}
		result[filepath] = fileFinalState
	}
	return result, nil
}

func BuildFileFromPatches(patches []domain.Diff) (string, error) {
	return BuildFileFromPatchesAndStartText("", patches)
}

func BuildFileFromPatchesAndStartText(startText string, patches []domain.Diff) (string, error) {
	dmp := diffmatchpatch.New()
	text := startText

	for i, patch := range patches {
		hunks := splitUnifiedIntoHunks(patch.PatchText)

		for j, hunk := range hunks {
			patches, err := dmp.PatchFromText(hunk)
			if err != nil {
				return "", fmt.Errorf("patch %d hunk %d: bad patch text: %w", i, j, err)
			}

			newText, applied := dmp.PatchApply(patches, text)
			for k, ok := range applied {
				if !ok {
					return "", fmt.Errorf("patch %d hunk %d: patch %d failed to apply",
						i, j, k)
				}
			}
			text = newText
		}
	}

	return text, nil
}

// splitUnifiedIntoHunks takes a unified-diff string that may contain
// multiple hunks and returns a []string, each starting with its "@@" header.
func splitUnifiedIntoHunks(patch string) []string {
	lines := strings.Split(patch, "\n")
	var hunks []string
	var current []string

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") && len(current) > 0 {
			// start of a new hunk: save the previous one
			hunks = append(hunks, strings.Join(current, "\n"))
			current = nil
		}
		current = append(current, line)
	}
	if len(current) > 0 {
		hunks = append(hunks, strings.Join(current, "\n"))
	}

	return hunks
}
