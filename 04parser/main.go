package main

import (
	"fmt"
	"04parser/parser"
)

func main() {
	links, err := parser.ParseLinks("examples/ex0.html")
	if err != nil {
		panic(err)
	}

	// Print results from ParseLinks function
	for _, v := range links {
		fmt.Println("---------- ----------")
		fmt.Printf("href:%v\ntext:%v \n", v.Href, v.Text)
	}
}