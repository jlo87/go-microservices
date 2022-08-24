package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jlo87/go-microservices/handlers"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {

	env.Parse()
	// logger
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// Create references to the hello and goodbye handlers
	ph := handlers.NewProducts(l)

	// Create a new servemux
	sm := http.NewServeMux()

	// Register the endpoints
	sm.Handle("/", ph)

	// Create a new server
	s := &http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		IdleTimeout:  120 * time.Second, // max time to for connections using TCP Keep-Alive
		ReadTimeout:  1 * time.Second,   // max time to read request from the cient
		WriteTimeout: 1 * time.Second,   // max time to write request to the client
	}
	// Handle ListenAndServe in a go func to prevent blocking
	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}

	}()

	// Broadcast message on this channel whenever an os.Kill/Interrupt is recieved
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Block here; reading from channel will block until there
	// is a message to be consumed
	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	// Wait until requests handled by the server have completed, then shutdown
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
