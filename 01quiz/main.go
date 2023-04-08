package main

import (
	"fmt"
	"encoding/csv"
	"io"
	"bufio"
	"strings"
	"os"
	"flag"
	"time"
)

var (
	// quiz's hits and misses counter
	hits, misses int	

	// Flag variables

	filename *string = flag.String("filename", "problems.csv", "choose quiz file inside quizzes folder")
	timer *int = flag.Int("timer", 30, "choose timer duration")
)

// !TODO: Add string trimming/cleaning (delete all blank spaces, no answer will ever have blank spaces)
func main() {	
	flag.Parse()

	// add time limit to quizz
	AddTimeLimit(*timer)

	// open csv to read it and defer its closing
	file := OpenFile("quizzes/" + *filename)
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


	fmt.Printf("\nQuiz ended, hits: %v, misses:%v\n", hits, misses)
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

// ReadCSVRecord reads a single line of a CSV file for each call (using *csv.Reader.Read()),
// it receives a pointer instead of a copy since Reader object has fields that may be
// modified by the client application.
// Returns nil if an io.EOF err is returned by reading a line of the csv, if
// another type of error is returned by reading a line, the function exits the program.
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

// AddTimeLimit adds a time limit to the quiz and ends it when the amount of time
// specified has elapsed, halting the current quiz's question and showing player's
// final stats.
func AddTimeLimit(amount int) {
	// AfterFunc triggers a callback after n seconds
	time.AfterFunc(time.Second * time.Duration(amount), func (){
		file := OpenFile("quizzes/" + *filename)
		defer file.Close()
		
		SkipNLines(file, hits + misses, func() {
			misses++
		})

		fmt.Printf("\nQuiz ended, hits: %v, misses:%v\n", hits, misses)
		os.Exit(0)
	})
}
// SkipNLines skips n lines from a file and calls a given function on all the
// non-skipped lines.
func SkipNLines(file *os.File, n int, fx func()) {
	csvReader := csv.NewReader(file)
	lineCount := 0
	for {
		// get a line's record and store it
		record := ReadCSVRecord(csvReader)
		if record == nil {
			return
		}

		// call the indicated funtion after n line skips
		lineCount++
		if lineCount > n {
			fx()
		}
	}
}