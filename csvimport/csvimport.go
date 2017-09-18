package csvimport

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type getCSV struct {
	path string
	rows map[int]map[string]string
}

// Converting CSV file to HashMap
func (c *getCSV) importFile(fp string) map[int]map[string]string {

	csvRows := make(map[int]map[string]string)
	var titles []string

	// Validating if option is CSV file and exists
	validFile, _ := regexp.MatchString(`\.csv$`, fp)
	if validFile {
		_, err := os.Stat(fp)
		if err != nil {
			fmt.Println("\nImport file does not exists:", fp)
			fmt.Println(err.Error() + "\n")
			os.Exit(1)
		}
	} else {
		fmt.Println("\nImport file must be in CSV:", fp+"\n")
		os.Exit(1)
	}

	f, _ := os.Open(fp)
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ','
	r.FieldsPerRecord = -1

	d, err := r.ReadAll()
	if err != nil {
		panic("Can't read csv file\n")
	}

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
	c.rows = csvRows
	return csvRows
}
