package main

import (
	"context"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"metrics_query/config"
	"metrics_query/events"
	"metrics_query/handlers"
	"metrics_query/service"
	"metrics_query/store"
	"metrics_query/stream"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger := log.New(os.Stdout, "[metrics-query] ", log.LstdFlags)

	//esdb client

	cfg := config.NewConfig()
	connString := fmt.Sprintf("esdb://%s:%s@%s:%s?tls=false", cfg.ESDBUser, cfg.ESDBPass, cfg.ESDBHost, cfg.ESDBPort)
	settings, err := esdb.ParseConnectionString(connString)
	if err != nil {
		log.Fatal(err)
	}
	esdbClient, err := esdb.NewClient(settings)
	if err != nil {
		log.Fatal(err)
	}
	eventStream, err := stream.NewESDBEventStream(esdbClient, cfg.ESDBGroup)
	if err != nil {
		log.Fatal(err)
	}

	mongoService, err := service.New(timeoutContext, logger)
	if err != nil {
		log.Fatal(err)
	}
	dbStore := store.NewStore(mongoService.GetCli())
	handler := events.NewEventHandler(dbStore)
	go eventStream.Process(handler.Handle)

	port := os.Getenv("PORT")
	router := mux.NewRouter()
	accommodationHandler := handlers.NewAccommodationHandler(dbStore)

	router.HandleFunc("/get/{id}/{period}", accommodationHandler.Get).Methods("GET")
	router.HandleFunc("/get/uuid", accommodationHandler.GenUUID).Methods("GET")

	headersOk := gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := gorillaHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	originsOk := gorillaHandlers.AllowedOrigins([]string{"http://localhost:4200"})
	server := http.Server{
		Addr:         ":" + port,
		Handler:      gorillaHandlers.CORS(headersOk, methodsOk, originsOk)(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	log.Println("Server listening on port", port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Panicf("PANIC FROM AUTH-SERVICE ON LISTENING")
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	log.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		log.Fatalf("Cannot gracefully shutdown...")
	}
	log.Println("Server stopped")

}
