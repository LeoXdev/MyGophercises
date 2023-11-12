package main

import (
	"fmt"
	"encoding/csv"
	"io"
	"bufio"
	"os"
	"flag"
	"time"
	"strings"
	"math/rand"
)

var (
	// quiz's hits and misses counter
	hits, misses int	

	// Flag variables

	filename *string = flag.String("filename", "problems.csv", "choose quiz file inside quizzes folder")
	timer *int = flag.Int("timer", 30, "choose timer duration")
	shuffle *bool = flag.Bool("shuffle", false, "shuffle quiz order?")
)

func main() {	
	flag.Parse()

	// add time limit to the quiz
	addTimeLimit(*timer)

	// shuffle if bool flag is true, panics if err occurs
	if *shuffle {
		shuffleCsv()
	}

	// open csv to read it and defer its closing
	file := openFile("quizzes/" + *filename)
	defer file.Close()
	csvReader := csv.NewReader(file)
	for {
		// a record is a single line of the quiz
		record := readCSVRecord(csvReader)
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
		// ReadString will read until the delimiter is entered
		// input includes the read delimiter
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		// remove the delimeter from the string
		input = strings.TrimSuffix(input, "\n")
		// remove white spaces
		input = cleanAnswer(input)

		// check if answer was right
		if input == record[1] {
			hits++
		} else {
			misses++
		}
	}

	fmt.Printf("\nQuiz ended, hits: %v, misses:%v\n", hits, misses)
	os.Exit(0)
}

// openFile opens a file and returns a pointer to it, doesn't return an error
// as the method itself manages it.
func openFile(filepath string) *os.File {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return file
}

// readCSVRecord reads a single line of a CSV file for each call (using *csv.Reader.Read()),
// it receives a pointer instead of a copy since Reader object has fields that may be
// modified by the client application.
// Returns nil if an io.EOF err is returned by reading a line of the csv, if
// another type of error is returned by reading a line, the function exits the program.
func readCSVRecord(r *csv.Reader) []string {
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

// addTimeLimit adds a time limit to the quiz and ends it when the amount of time
// specified has elapsed, halting the current quiz's question and showing player's
// final stats.
func addTimeLimit(amount int) {
	// AfterFunc triggers a callback after n seconds
	time.AfterFunc(time.Second * time.Duration(amount), func (){
		file := openFile("quizzes/" + *filename)
		defer file.Close()
		
		// sum of hits + misses equals to the total answered questions
		skipNLines(file, hits + misses, func() {
			misses++
		})

		fmt.Printf("\nQuiz ended, hits: %v, misses:%v\n", hits, misses)
		os.Exit(0)
	})
}
// skipNLines skips n lines from a file and calls a given function on all the
// non-skipped lines.
func skipNLines(file *os.File, n int, fx func()) {
	csvReader := csv.NewReader(file)
	lineCount := 0
	for {
		// get a line's record and store it
		record := readCSVRecord(csvReader)
		if record == nil {
			return
		}

		// call the indicated function after n line skips
		lineCount++
		if lineCount > n {
			fx()
		}
	}
}
// cleanAnswer removes white spaces inside an entered answer, the strings
// wouldn't need anymore processing.
func cleanAnswer(str string) string {
	return strings.ReplaceAll(str, " ", "")
}
// shuffleCsv shuffles the quiz using the modern version of Fisher-Yates
// algorithm.
func shuffleCsv() {
	file := openFile("quizzes/" + *filename)
	defer file.Close()
	csvReader := csv.NewReader(file)

	csvData, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(csvData)
	// Fisher-Yates algorithm for shuffling lines of the csvData
	n := len(csvData) - 1
	rand.Seed(time.Now().UnixNano())
	for i := n; i > 0; i-- {
		// select a random index from the unshuffled part of the csvData
		// adding 1 to i leaves the possibility of an reciprocal exchange
		j := rand.Intn(i + 1)
		csvData[i], csvData[j] = csvData[j], csvData[i]
	}
	fmt.Println(csvData)

	// os.Create will truncate selected quiz file
	file, err = os.Create("quizzes/" + *filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// attempt to write shuffled lines into file
	err = writer.WriteAll(csvData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}