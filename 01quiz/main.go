package main

import (
	"encoding/csv"
	"fmt"
	"strings"
	"io"
	"bufio"
	"os"
)

func main() {
	var (
		// quiz's hits and misses counter
		hits, misses int = 0, 0
	)

	// open csv to read it and defer its closing
	file, err := os.Open("problems.csv")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	for {
		record, err := csvReader.Read()
		// double error check as Read() may return different errors
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// csv records processing,
		// first field is always a question
		// second field is always the answewr to the question
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

func ReadFile() {

}
func ReadCSVRecord(r csv.Reader) {

}