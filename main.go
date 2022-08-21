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
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "Oops", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(rw, "Hello %s", d)
	})

	// When the path matches goodbye, it executes this function
	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("Goodbye World")
	})

	http.ListenAndServe(":9090", nil)
}
