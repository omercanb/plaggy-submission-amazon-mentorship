package domain

import "time"

type Diff struct {
	ID        uint
	FilePath  string
	PatchText string
	Timestamp time.Time
}
