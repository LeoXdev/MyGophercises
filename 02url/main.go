package main

import (
	"fmt"
	"net/http"

	"02url/handler"
)

func main() {
	// create fallback handler
	mux := http.DefaultServeMux
	mux.HandleFunc("/hello", hello)

	pathsToUrls := map[string]string {
		"/go" : "https://go.dev/",
	}
	
	// open server with both paths to redirect and fallback handler
	http.ListenAndServe(":8080", handler.MapHandler(pathsToUrls, mux))
}

// hello contains the handler for the '/hello' path (fallback handler)
func hello(w http.ResponseWriter,r *http.Request) {
	fmt.Fprintln(w, "fallback handler...")
}
