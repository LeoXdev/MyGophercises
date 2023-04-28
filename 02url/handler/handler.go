package handler

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// creation of a HandleFunc to have access to request's URL
	return func(w http.ResponseWriter, r *http.Request) {		
		// k holds paths to be redirected, v their respective endpoints (of going to k)
		for k, v := range pathsToUrls {
			// if path in the request applies to be redirected (is present on our map),
			// it gets redirected to the linked endpoint (v)
			// else, the fallback handler is executed
			if k == r.URL.Path {
				http.Redirect(w, r, v, http.StatusSeeOther)
			} else {
				fallback.ServeHTTP(w, r)
			}
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(string(yaml))
	if err != nil {
	  return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}
// --- Helpher Functions ---

// pathToURL is a struct for generating YAML that'll hold an attribute for
// the path to have a redirection and another one that indicates the endpoint of
// the redirection.
type pathToURL struct {
	Path string	`yaml:"path"`
	URL string	`yaml:"url"`
}
// parseYAML parses a string into a []byte,
// parse stands for read and try to translate.
func parseYAML(yamlToParse string) ([]byte, error) {
	var parsedYAML pathToURL

	// parsingErr will only hold errors caused by a bad YAML syntax
	parsingErr := yaml.Unmarshal([]byte(yamlToParse), &parsedYAML)

	a, err := yaml.Marshal(&parsedYAML)
	if err != nil {
		panic(err)
	}
	
	return a, parsingErr
}
// buildMap writes a yaml's values into a map and returns it.
func buildMap(yamlData []byte) map[string]string {
	mapWithValues := make(map[string]string)

	yaml.Unmarshal(yamlData, &mapWithValues)

	mapToReturn := make(map[string]string)
	var path, url string
	for k, v := range mapWithValues {
		if k == "path" {
			path = v
		} else if k == "url" {
			url = v
		}
	}
	mapToReturn[path] = url
	return mapToReturn
}
