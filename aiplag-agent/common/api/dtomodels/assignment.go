package dtomodels

import "time"

type Assignment struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	AssignedAt time.Time `json:"assignedAt"`
	DueDate    time.Time `json:"dueDate"`
}
