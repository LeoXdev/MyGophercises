package main

import (
	"net/http"

	"03adventure/handler"
)

// pagina home para escoger historia
// intro de la historia seleccionada y options al final
func main() {
	indexHandler := handler.NewIndexHandler()
	gopherJsonHandler := handler.NewGopherJsonHandler()

	http.Handle("/index", indexHandler)
	http.Handle("/adventures/gopher", gopherJsonHandler)

	// Note about handler identifier:
	// Code editor may include identifiers from 02url/handler but they cannot be used here
	
	http.ListenAndServe(":8080", nil)
}
