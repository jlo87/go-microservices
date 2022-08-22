package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jlo87/go-microservices/handlers"
)

func main() {
	// logger
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// Create references to the hello and goodbye handlers
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	// Create a new servemux
	sm := http.NewServeMux()

	// Register the endpoints
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	// Handle ListenAndServe in a go func to prevent blocking
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// Broadcase message on this channel whenever an os.Kill/Interrupt is recieved
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Block here; reading from channel will block until there
	// is a message to be consumed
	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// Wait until requests handled by the server have completed, then shutdown
	s.Shutdown(tc)
}
