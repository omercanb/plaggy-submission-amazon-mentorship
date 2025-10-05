package models

import "time"

type Detection struct {
	ID         string    `json:"id"`
	CreatedBy  string    `json:"createdBy"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	Severity   uint      `json:"severity"`
	HomeworkID string    `json:"homeworkId"`
	FilePath   string    `json:"filePath"`
	DiffData   string    `json:"diffData"`
}
