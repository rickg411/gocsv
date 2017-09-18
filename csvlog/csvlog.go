package csvlog

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type csvGo struct {
	logFile *os.File
	count   int
	quit    chan bool
	row     chan []string
}

func (c csvGo) create(fileName string) (*os.File, error) {

	fp := []string{os.Getenv("HOME"), "Desktop", strings.Replace(os.Args[0], "./", "", -1)}
	path := strings.Join(fp, "/")

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0775)
	}

	if fileName != "" {
		fp = append(fp, fileName+"-")
	}

	fp = append(fp, time.Now().Format(time.RFC1123)+".csv")
	f, err := os.Create(strings.Join(fp, "/"))

	return f, err
}

//*** Log Functions ***//
func (c *csvGo) startLog(fileName string, title ...[]string) *csvGo {

	c = &csvGo{row: make(chan []string)}

	f, err := c.create(fileName)
	if err != nil {
		fmt.Println("Could not create logfile. This will prevent logging to:", fileName)
		recover()
	}
	c.logFile = f

	go c.initLogger()

	if len(title[0]) > 0 {
		c.addRow(title[0])
	}
	return c
}

func (c csvGo) addRow(row []string) {
	go func() {
		c.row <- row
	}()
}

func (c *csvGo) initLogger() {
	w := csv.NewWriter(c.logFile)
	w.UseCRLF = false
	defer c.logFile.Close()
	for {
		select {
		case row := <-c.row:
			c.count++
			w.Write(row)
			w.Flush()
		case <-c.quit:
			close(c.row)
		}
	}
}

func (c csvGo) close() int {
	close(c.row)
	return c.count
}
