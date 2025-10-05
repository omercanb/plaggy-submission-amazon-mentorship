// Computing diffs between two files
package filesystemwatching

import (
	"aiplag-agent/daemon/models"
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type FileDiffer struct {
	differ *diffmatchpatch.DiffMatchPatch
}

func NewFileDiffer() *FileDiffer {
	fd := &FileDiffer{}
	fd.differ = diffmatchpatch.New()
	return fd
}

func (fd *FileDiffer) Diff(file1 models.File, file2 models.File) (string, error) {
	// diffmatchpath PatchMake can sometimes throw a fatal error, especially when the file is a binary.
	// For these cases, we need to catch the error with recover()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("FileDiffer panic on %s vs %s: %v", file1.Path(), file2.Path(), r)
		}
	}()

	filecontent1, err := file1.Read()
	if err != nil {
		log.Printf("FileDiffer Diff: Failed to read %s: %v", file1.Path(), err)
		return "", err
	}
	filecontent2, err := file2.Read()
	if err != nil {
		log.Printf("FileDiffer Diff; Failed to read %s: %v", file2.Path(), err)
		return "", err
	}

	patchText := fd.UnifiedLineLevelPatches(filecontent1, filecontent2)
	fmt.Println(patchText)
	return patchText, nil
}

func (fd *FileDiffer) UnifiedLineLevelPatches(filecontent1 string, filecontent2 string) string {
	text1Lines, text2Lines, lineArray := fd.differ.DiffLinesToRunes(filecontent1, filecontent2)
	diffs := fd.differ.DiffMainRunes(text1Lines, text2Lines, false)
	diffs = fd.differ.DiffCharsToLines(diffs, lineArray)
	diffs = fd.differ.DiffCleanupSemantic(diffs)
	diffs = fd.differ.DiffCleanupEfficiency(diffs)
	patches := fd.differ.PatchMake(filecontent1, diffs)
	unfilteredPatchText := fd.differ.PatchToText(patches)

	// The next step is for removing the weird url query escaping of the diffmatchpatch library
	sb := strings.Builder{}
	for patchLine := range strings.SplitSeq(unfilteredPatchText, "\n") {
		if len(patchLine) == 0 {
			sb.WriteString("\n")
			continue
		}
		diffTypeChars := []byte{'+', '-', ' '}
		patchType := patchLine[0]
		encodedPatchText := patchLine[1:]

		i := slices.Index(diffTypeChars, patchType)
		if i == -1 {
			sb.WriteString(patchLine)
			sb.WriteString("\n")
			continue
		}

		patchHunk, err := url.QueryUnescape(encodedPatchText)
		if err != nil {
			sb.WriteString(patchLine)
			sb.WriteString("\n")
			continue
		}

		for hunkLine := range strings.SplitSeq(patchHunk, "\n") {
			sb.WriteByte(patchType)
			sb.WriteString(hunkLine)
			sb.WriteString("\n")
		}
	}
	return sb.String()
	/*
			text1, text2, lineArray := fd.differ.DiffLinesToChars(filecontent1, filecontent2)
			diffs := fd.differ.DiffMain(text1, text2, false)
			lineDiffs := fd.differ.DiffCharsToLines(diffs, lineArray)

			var sb strings.Builder
			patches := fd.differ.PatchMake(filecontent1, filecontent2)
			patchTextForHeader := fd.differ.PatchToText(patches)
			headerEndIdx := strings.Index(patchTextForHeader, "\n")
			if headerEndIdx != -1 {
				header := patchTextForHeader[:headerEndIdx]
				sb.WriteString(header)
				sb.WriteString("\n")
			}

			for _, diff := range lineDiffs {
				for line := range strings.SplitSeq(diff.Text, "\n") {
					switch diff.Type {
					case diffmatchpatch.DiffInsert:
						sb.WriteString("+ ")
					case diffmatchpatch.DiffDelete:
						sb.WriteString("- ")
					case diffmatchpatch.DiffEqual:
						sb.WriteString("  ")
					}

					// Trim trailing newline if present
					line = strings.TrimRight(line, "\n")
					sb.WriteString(line)
					sb.WriteByte('\n')
				}
			}
		return sb.String()
	*/
}
