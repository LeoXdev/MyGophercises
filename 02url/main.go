// The server may not start if port :8080 still has a process running, to kill it
// enter 'fuser -k 8080/tcp' in terminal
package main

import (
	"flag"
	"fmt"
	"net/http"

	"io/ioutil"

	"02url/handler"
)

var (
	useFile *bool = flag.Bool("usefile", false, "use file at path datafiles/data.yaml ?")
)

func main() {
	flag.Parse()

	// create fallback handler
	mux := http.DefaultServeMux
	mux.HandleFunc("/hello", hello)

	// ------ MapHandler implementation ------
	//pathsToUrls := map[string]string {
	//	"/go" : "https://go.dev/",
	//}
	
	// open server with both paths, to redirect and fallback handler
	//http.ListenAndServe(":8080", handler.MapHandler(pathsToUrls, mux))

	// ------ YAMLhandler implementation ------
	var data []byte
	if (*useFile) {
		yamlFile, err := ioutil.ReadFile("datafiles/data.yaml")
		if err != nil {
			panic(err)
		}

		data = yamlFile
	} else {
		data = []byte(`
path: /google
url: https://www.google.com/
`)
	}
	yamlHandler, err := handler.YAMLHandler(data, mux)
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", yamlHandler)
}

// hello contains the handler for the '/hello' path (fallback handler)
func hello(w http.ResponseWriter,r *http.Request) {
	fmt.Fprintln(w, "fallback handler...")
}
