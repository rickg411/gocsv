package gocsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// GoCsv - Adding comment here
type GoCsv struct {
	file   *os.File
	writer *csv.Writer
	count  int
	quit   chan bool
	row    chan []string
}

//*** Log Functions ***//

// Start - Initalize Receiver and creating csv file
func Start(fileName string, title ...[]string) *GoCsv {

	c := new(GoCsv)
	c = &GoCsv{row: make(chan []string)}

	f, err := c.createLog(fileName)
	if err != nil {
		fmt.Println("Could not create logfile. This will prevent logging to:", fileName)
		recover()
	}

	c.file = f
	c.writer = csv.NewWriter(f)
	c.writer.UseCRLF = false

	go c.logger()

	if len(title[0]) > 0 {
		c.AddRow(title[0])
	}
	return c
}

// Creates the csv file
func (c GoCsv) createLog(name string) (*os.File, error) {

	fp := []string{os.Getenv("HOME"), "Desktop", strings.Replace(os.Args[0], "./", "", -1)}
	path := strings.Join(fp, "/")

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0775)
	}

	if name != "" {
		fp = append(fp, name+" ")
	}
	// creating unquie file name not ot overwrite existing
	// Need to revamp later
	fp = append(fp, time.Now().Format(time.RFC1123)+".csv")
	f, err := os.Create(strings.Join(fp, "/"))

	return f, err
}

// AddRow - Adding csv row
func (c GoCsv) AddRow(row []string) {
	go func() {
		c.row <- row
	}()
}

// Channel reciever to write to file
func (c *GoCsv) logger() {
	for {
		select {
		case row := <-c.row:
			c.count++
			c.writer.Write(row)
			c.writer.Flush()
		}
	}
}

// Close - Closing the channel and file
func (c GoCsv) Close() int {
	c.file.Close()
	close(c.row)
	return c.count
}
