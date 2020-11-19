package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func initDatabase(pgHost, pgPort, pgUser, pgPassword, dbName string) *sqlx.DB {
	log.Println("Connecting to Database:", pgHost)

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost,
		pgPort,
		pgUser,
		pgPassword,
		dbName,
	))
	if err != nil {
		log.Fatalf("Error: connecting to Database %v", err)
	}

	log.Println("Connected to Database")
	return db
}

func main() {
	// envs
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error: loading env %v", err)
		return
	}
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// init DB
	db := initDatabase(pgHost, pgPort, pgUser, pgPassword, dbName)
	defer db.Close()

	// init cache
	cachePool := initCache()
	cacheApi, err := NewCacheAPI(cachePool)
	if err != nil {
		log.Printf("Error: initializing cacheApi %v", err)
		return
	}
	// create router
	router := mux.NewRouter()

	// create new app instance
	_, err = New(router, db, &cacheApi)
	if err != nil {
		log.Printf("Error: initializing app %v", err)
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

	log.Println("Server is running on port 3000...")

	// blocking: empty the channel so it can unblock
	<-gracefulShutdownChan

	// start graceful server shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server graceful shutdown failed: %+s", err)
		cancel()
	}

	log.Println("Server Stopped")
	os.Exit(0)
}
