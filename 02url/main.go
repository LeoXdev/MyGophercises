package main

import (
	"fmt"
	"net/http"

	"02url/handler"
)

func main() {
	var yamlData string = `
path: /google
url: https://www.google.com/
`

	// create fallback handler
	mux := http.DefaultServeMux
	mux.HandleFunc("/hello", hello)

	// MapHandler implementation
	//pathsToUrls := map[string]string {
	//	"/go" : "https://go.dev/",
	//}
	
	// open server with both paths, to redirect and fallback handler
	//http.ListenAndServe(":8080", handler.MapHandler(pathsToUrls, mux))

	// YAMLhandler implementation
	yamlHandler, err := handler.YAMLHandler([]byte(yamlData), mux)
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", yamlHandler)
}

// hello contains the handler for the '/hello' path (fallback handler)
func hello(w http.ResponseWriter,r *http.Request) {
	fmt.Fprintln(w, "fallback handler...")
}
