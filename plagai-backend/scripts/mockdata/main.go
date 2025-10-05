package main

import "github.com/plagai/plagai-backend/scripts/mockdata/mockdiff"

type PatchWithTimestamp struct {
	Patch     string `json:"patch_text"`
	Timestamp int64  `json:"timestamp"`
}

/*
mockdata is a directory for mocking the data that goes into our backend like diffs and assignment submissions
Currently it supports creating mock diffs from a file
In main uncomment any creation function to enable a cli to create the mock data
*/
func main() {
	mockdiff.Start()
}
