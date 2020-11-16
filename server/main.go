package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// envs

	router := mux.NewRouter()

	// create new app instance
	_, err := New(router)
	if err != nil {
		log.Println(err)
		return
	}

	// graceful shutdown
	log.Println("Starting Server...")
	gracefulShutdownChan := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdownChan, os.Interrupt, os.Kill)

	// start server routine
	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
			close(gracefulShutdownChan)
		}
	}()

	log.Println("Server is running...")

	// blocking: empty the channel so it can unblock
	<-gracefulShutdownChan

	// start graceful server shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server graceful shutdown failed: %+s", err)
		cancel()
	}

	log.Println("Stopped Server")
	os.Exit(0)
}
