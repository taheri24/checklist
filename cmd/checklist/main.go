package main

import (
	"flag"
	"fmt"
	"os"

	"checklist/internal/checklist"
)

func main() {
	checklistPath := flag.String("input", "checklist.txt", "path to checklist file with one item per line")
	outputPath := flag.String("output", "selected.txt", "path where selected items will be written")
	flag.Parse()

	if err := checklist.Run(*checklistPath, *outputPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
