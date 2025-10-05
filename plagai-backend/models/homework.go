package models

import "time"

type Homework struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	AssignedAt time.Time `json:"assignedAt"`
	DueDate    time.Time `json:"dueDate"`
}
