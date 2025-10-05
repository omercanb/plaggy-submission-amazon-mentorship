package mockdiff

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/plagai/plagai-backend/core"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type MockFileState struct {
	Text string
	Time time.Time
}

type PatchRecord struct {
	Patch     string `json:"patch_text"`
	Timestamp int64  `json:"timestamp"`
}

// Creates a CLI thats used for generating a JSON of mock diffs and timestamps for a file
// The diffs are created as if the files was written top to bottom, with 5-10 lines added each diff
func Start() {
	filePtr := flag.String("file", "", "the file to create mock diffs for")
	flag.Parse()

	if *filePtr == "" {
		fmt.Println("Error: -file argument is required")
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.Open(*filePtr)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	patchRecords := simulateCreationOfFile(file)

	jsonBytes, err := json.MarshalIndent(patchRecords, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal failed: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func simulateCreationOfFile(file *os.File) []PatchRecord {
	mockChanges := splitLinesRandomly(file, core.Range{Min: 5, Max: 10})
	mockPatches := convertToPatches(mockChanges)

	mockFileEditStartingTime := time.Now().Add(-48 * time.Hour)
	simulatedTime := mockFileEditStartingTime

	patchRecords := []PatchRecord{}
	for _, patch := range mockPatches {
		timeTakenForChange := time.Duration(len(patch)*10) * time.Minute
		simulatedTime = simulatedTime.Add(timeTakenForChange)

		patchRecords = append(patchRecords, PatchRecord{
			Patch:     patch,
			Timestamp: simulatedTime.UnixMilli(),
		})

	}
	return patchRecords
}

func convertToPatches(changes []string) []string {
	accumulator := ""
	patches := []string{}
	dmp := diffmatchpatch.New()

	for _, change := range changes {
		patch := dmp.PatchMake(accumulator, change)
		newAccumulator, applied := dmp.PatchApply(patch, accumulator)
		accumulator = newAccumulator
		validatePatches(applied, patch)

		patches = append(patches, dmp.PatchToText(patch))
	}

	return patches
}

func validatePatches(applications []bool, patch []diffmatchpatch.Patch) {
	for _, application := range applications {
		if !application {
			log.Fatalf("Patch application failed in patch: %s", patch)
		}
	}
}

func splitLinesRandomly(file *os.File, lineRange core.Range) []string {
	lines := core.ReadLines(file)
	mockChanges := []string{}
	currentChangeAccumulator := []string{}
	linesToAddToCurrentChange := rand.IntN(lineRange.Max-lineRange.Min+1) + lineRange.Min

	for _, line := range lines {
		currentChangeAccumulator = append(currentChangeAccumulator, line)
		linesToAddToCurrentChange -= 1

		if linesToAddToCurrentChange <= 0 {
			currentChangeString := strings.Join(currentChangeAccumulator, "")
			mockChanges = append(mockChanges, currentChangeString)
			currentChangeAccumulator = []string{}

			linesToAddToCurrentChange = rand.IntN(lineRange.Max-lineRange.Min+1) + lineRange.Min
		}
	}

	if len(currentChangeAccumulator) > 0 {
		currentChangeString := strings.Join(currentChangeAccumulator, "")
		mockChanges = append(mockChanges, currentChangeString)
	}

	return mockChanges
}
