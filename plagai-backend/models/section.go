package models

type Section struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	HomeworkIDs string `json:"homeworkIds"`
}
