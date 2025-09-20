package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	// Replace with your actual CSV file path
	filePath := "tags.csv"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header first
	header, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read header: %v", err)
	}

	// Find the index of "tag_name" column
	tagIndex := -1
	for i, h := range header {
		if h == "tag_name" {
			tagIndex = i
			break
		}
	}
	if tagIndex == -1 {
		log.Fatalf("tag_name column not found in header")
	}

	maxLen := 0
	rowCount := 0
	longestTagName := ""

	// Iterate through rows
	for {
		record, err := reader.Read()
		if err != nil {
			// Stop at EOF
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("Error reading record: %v", err)
		}

		rowCount++
		tag := record[tagIndex]
		if len(tag) > maxLen {
			maxLen = len(tag)
			longestTagName = tag
		}
	}

	fmt.Printf("Processed %d rows\n", rowCount)
	fmt.Printf("Longest tag_name length: %d\n", maxLen)
	fmt.Printf("Longest tag_name: %s\n", longestTagName)
}
