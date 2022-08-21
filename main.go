// Writing a basic web server in Go

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Basic HTTP Handler
	// HandleFunc is a convenience method on the go HTTP package
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Hello World")

		// Read everything from the body into the variable d, the data
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "Oops", http.StatusBadRequest)
			return
		}
		// Print string back to the rw which will be returned to the user
		fmt.Fprintf(rw, "Hello %s", d)
	})

	// When the path matches goodbye, it executes this function
	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("Goodbye World")
	})

	// Convenience method that constructs an HTTP server and registers the default handler
	http.ListenAndServe(":9090", nil)
}
