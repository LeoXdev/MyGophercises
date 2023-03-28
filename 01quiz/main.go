package main

import (
	"encoding/csv"
	"fmt"
	"strings"
	"io"
	"bufio"
	"os"
	"flag"
)

func main() {
	var (
		// quiz's hits and misses counter
		hits, misses int = 0, 0

		// Flag variables
		filename *string = flag.String("filename", "problems.csv", "choose quiz file inside quizzes folder")
	)
	flag.Parse()

	// open csv to read it and defer its closing
	file := OpenFile("quizzes/" + *filename) // !TODO: change the filename literal for a flag variable
	defer file.Close()

	csvReader := csv.NewReader(file)
	for {
		record := ReadCSVRecord(csvReader)
		// in case no more records are left
		if record == nil {
			break
		}

		// csv records processing,
		// first field is always a question
		// second field is always the answer to the question
		fmt.Printf("%v", record[0])
		fmt.Print("=? ")
		reader := bufio.NewReader(os.Stdin)
		// ReadString will block until the delimiter is entered
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		// remove the delimeter from the string, ReadString returns a value up to the ocurring parameter
		input = strings.TrimSuffix(input, "\n")

		if input == record[1] {
			hits++
		} else {
			misses++
		}
	}


	fmt.Printf("Quiz ended, hits: %v, misses:%v\n", hits, misses)
	os.Exit(0)
}

// OpenFile opens a file and returns a pointer to it, doesn't return an error
// as the method itself manages it.
func OpenFile(filepath string) *os.File {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return file
}

// ReadCSVRecord reads a single line of a CSV file for each call, it receives
// a pointer instead of a copy since Reader object has fields that may be
// modified by the client application.
func ReadCSVRecord(r *csv.Reader) []string {
	record, err := r.Read()

	// double error check as Read() may return different errors
	if err == io.EOF {
		// without this check, an error will be thrown and the program will
		// shut down when reaching the end of the CSV
		return nil
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return record
}