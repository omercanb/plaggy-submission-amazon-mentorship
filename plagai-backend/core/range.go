package core

import "log"

type Range struct {
	Min int
	Max int
}

func NewRange(min, max int) Range {
	if min > max {
		log.Fatalf("Invalid range %d, %d\n", min, max)
	}
	return Range{Min: min, Max: max}
}
