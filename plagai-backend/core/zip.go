package core

import "fmt"

func Zip2[A, B any](a []A, b []B) []Pair[A, B] {
	if len(a) != len(b) {
		panic(fmt.Sprintf("Zip2 error: slice length mismatch %d vs %d", len(a), len(b)))
	}

	size := len(a)
	result := make([]Pair[A, B], size)

	for i := range a {
		result = append(result, Pair[A, B]{a[i], b[i]})
	}

	return result
}
