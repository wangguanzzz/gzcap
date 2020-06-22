package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var filename string
var filecontext [][]string
var filenum int

func main() {
	// open file
	fs, err := os.Open("luowen.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer fs.Close()

	r := csv.NewReader(fs)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			fmt.Println("reading error")
		}
		tmp := strings.ReplaceAll(row[0], "/", "-")
		datefile := fmt.Sprintf("%s.csv", tmp)
		if datefile != filename {
			//end wrting file
			writeCSV(filename, filecontext)
			filename = datefile
			filecontext = [][]string{{}}
			filenum++
			fmt.Println("file add ", filenum)
		}

		filecontext = append(filecontext, row)

	}
}

func writeCSV(filename string, context [][]string) {
	if filename == "" {
		return
	}
	file, err := os.Create(filename)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range context {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
