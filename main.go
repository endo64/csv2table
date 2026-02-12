package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"io"

	"github.com/olekukonko/tablewriter"
)

func main() {
	// Parse arguments
	setRowLine := flag.Bool("rowline", true, "Draws a line between each row")
	header := flag.Bool("header", true, "Use the first row as the header")
	cols := flag.Int("cols", 0, "Maximum number of columns")
	flag.Parse()

	// Check for the CSV file path argument
	if len(flag.Args()) == 0 {
		fmt.Println("Usage: csv2table [rowline|header|cols] [csv file]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var content [][]string = make([][]string, 0, 0)

	// Open the CSV file
	file, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	// Read the CSV file
	csvFile := csv.NewReader(file)

	if *cols <= 0 {
		content, err = csvFile.ReadAll()
		if err != nil {
			log.Fatalf("Error reading CSV file: %v", err)
		}
	} else {
		// Read the file record by record
		for {
			record, err := csvFile.Read()
			if err == io.EOF {
				break // End of File
			}
			if err != nil {
				log.Fatalf("Error reading CSV file: %v", err)
			}

			// Process only up to cols
			data := make([]string, 0, *cols)
			for i, field := range record {
				if i >= *cols {
					break // Stop processing once the column limit is reached
				}
				data = append(data, field)
			}

			// Use the data (which has at most cols fields)
			content = append(content, data)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)

	// Should the first row be used as the table header?
	firstRow := 0
	if *header {
		table.SetHeader(content[0])
		firstRow = 1
	}

	// Keep all columns to a single line
	table.SetAutoWrapText(false)

	// Add a line delineating each row
	table.SetRowLine(*setRowLine)

	// Append the rest of the CSV rows to the table
	for _, row := range content[firstRow:] {
		table.Append(row)
	}

	// Render the table to the console
	table.Render()
}
