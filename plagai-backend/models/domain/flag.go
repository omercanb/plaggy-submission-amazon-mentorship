package domain

type Flag struct {
	ID              uint
	Diff            Diff
	FlagExplanation string
	Severity        int
}
