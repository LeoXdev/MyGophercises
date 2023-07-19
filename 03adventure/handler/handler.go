package handler

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

// IndexHandler implements the http.Handler type/interface.
type IndexHandler struct{}
func NewIndexHandler() *IndexHandler {
	// We don't initialize any field because we're only using methods for this struct
	instance := &IndexHandler{}
	return instance
}
// ServeHTTP implements the method required by the interface http.Handler.
// Registers the index page in the ServeMux which is used to select an adventure and start it.
func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Your custom logic goes here
	// Use w to write response to client and r to read request from client
	// Read the template string from file
	tpl, _ := ioutil.ReadFile("templates/index.html")

	// Parse the template
	t := template.Must(template.New("tpl").Parse(string(tpl)))

	// Provide data
	data := struct {
		Title   string
		Header  string
		Content string
	}{
		Title:   "Adventure",
		Header:  "Exercise 03: Choose your own adventure",
	}

	// Write template with values in w
	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GopherJsonHandler implements the http.Handler type/interface.
type GopherJsonHandler struct{}
// GopherJsonStruct holds the structure required to hold the gopher json in memory during runtime.
// Obtained trough: https://transform.tools/json-to-go
type GopherJsonStruct struct {
	Intro struct {
		Title   string   `json:"title"`
		Story   []string `json:"story"`
		Options []struct {
			Text string `json:"text"`
			Arc  string `json:"arc"`
		} `json:"options"`
	} `json:"intro"`
	NewYork struct {
		Title   string   `json:"title"`
		Story   []string `json:"story"`
		Options []struct {
			Text string `json:"text"`
			Arc  string `json:"arc"`
		} `json:"options"`
	} `json:"new-york"`
	Debate struct {
		Title   string   `json:"title"`
		Story   []string `json:"story"`
		Options []struct {
			Text string `json:"text"`
			Arc  string `json:"arc"`
		} `json:"options"`
	} `json:"debate"`
	SeanKelly struct {
		Title   string   `json:"title"`
		Story   []string `json:"story"`
		Options []struct {
			Text string `json:"text"`
			Arc  string `json:"arc"`
		} `json:"options"`
	} `json:"sean-kelly"`
	MarkBates struct {
		Title   string   `json:"title"`
		Story   []string `json:"story"`
		Options []struct {
			Text string `json:"text"`
			Arc  string `json:"arc"`
		} `json:"options"`
	} `json:"mark-bates"`
	Denver struct {
		Title   string   `json:"title"`
		Story   []string `json:"story"`
		Options []struct {
			Text string `json:"text"`
			Arc  string `json:"arc"`
		} `json:"options"`
	} `json:"denver"`
	Home struct {
		Title   string        `json:"title"`
		Story   []string      `json:"story"`
		Options []interface{} `json:"options"`
	} `json:"home"`
}
func NewGopherJsonHandler() *GopherJsonHandler {
	instance := &GopherJsonHandler{}
	return instance
}
// ServeHTTP implements the method required by the interface http.Handler.
// Registers the gopher json in the ServeMux, so it can be used by javascript to be rendered.
func (g *GopherJsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//current_path, _ := os.Getwd() MyGophercises/03adventure

	// Open and read the JSON file
	file, err := os.Open("adventures/gopher.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Decode the JSON data into a struct
	var data GopherJsonStruct
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the struct as JSON and write it to the response
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
