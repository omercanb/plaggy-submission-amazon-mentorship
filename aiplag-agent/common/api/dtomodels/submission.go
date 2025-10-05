package dtomodels

type Submission struct {
	AssignmentId uint        `json:"assignmentID"`
	Edits        []EditEvent `json:"edits"`
}
