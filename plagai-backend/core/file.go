package core

import (
	"bufio"
	"log"
	"os"
)

func ReadLines(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()+"\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}
