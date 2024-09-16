package main

import (
	// import go-chi v5
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	// create a channel to listen for shutdown signals
	shutdown := make(chan bool)
	// intercept signals
	signalsChan := make(chan os.Signal, 1)

	// listen for SIGINT and SIGTERM
	signal.Notify(signalsChan, syscall.SIGINT, syscall.SIGTERM)

	r := chi.NewRouter()
	r.Get("/liveness", liveness)
	r.Get("/readiness", readiness)
	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go server.ListenAndServe()

	// wait for a signal
	go func() {
		s := <-signalsChan
		fmt.Printf("signal received: %v\n", s)
		// shutdown the server before exiting
		shutdownCTX, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCTX)
		fmt.Println("server shutdown")
		shutdown <- true
	}()

	// wait for shutdown
	<-shutdown
}

func liveness(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func readiness(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
