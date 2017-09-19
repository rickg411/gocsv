package gocsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// GetCSV - Getting the csv file
func GetCSV(fp string) *os.File {

	if validFile, _ := regexp.MatchString(`\.csv$`, fp); !validFile {
		fmt.Println("\nExpected to get file with .csv extension\nInstead got:", fp+"\n")
		os.Exit(1)
	}

	f, err := os.Open(fp)
	if err != nil {
		panic(err.Error())
	}

	return f
}

func parseCSV(f *os.File) [][]string {
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ','
	r.FieldsPerRecord = -1

	d, _ := r.ReadAll()
	return d
}

// GetMap - Converting CSV rows to Map
func GetMap(fp string) map[int]map[string]string {

	csvRows := make(map[int]map[string]string)
	var titles []string

	d := parseCSV(GetCSV(fp))

	for i := 0; i < len(d[0]); i++ {
		titles = append(titles, strings.ToLower(d[0][i]))
	}

	// Converting CSV file to map
	for i, val := range d {
		// Skipping first row/ title row of csv
		if i > 0 {
			csvRow := make(map[string]string)
			for t, title := range titles {
				csvRow[title] = val[t]
			}
			csvRows[i] = csvRow
		}
	}

	return csvRows
}
