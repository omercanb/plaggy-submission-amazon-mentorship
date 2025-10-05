package models

import "time"

type Assignment struct {
	ID         uint
	Title      string
	DueDate    time.Time
	AssignedAt time.Time
}
