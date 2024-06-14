package main

import (
	// import go-chi v5
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
		server.Shutdown(nil)
		fmt.Println("server shutdown")
		shutdown <- true
	}()

	// wait for shutdown
	<-shutdown
}
