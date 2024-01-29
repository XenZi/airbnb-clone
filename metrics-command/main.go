package main

import (
	"context"
	"fmt"
	"log"
	"metrics-command/commands/handler"
	"metrics-command/config"
	"metrics-command/handlers"
	"metrics-command/store"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/EventStore/EventStore-Client-Go/esdb"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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

	eventStore := store.NewESDBStore(esdbClient)
	commandHandler := handler.NewHandler(eventStore)
	userHandler := handlers.NewUserHandler(commandHandler)
	reservationHandler := handlers.NewReservationHandler(commandHandler)
	_ = handlers.NewRatingHandler(commandHandler)

	port := os.Getenv("PORT")

	router := mux.NewRouter()
	router.HandleFunc("/joinedAt", userHandler.CreateJoinedAt).Methods("POST")
	router.HandleFunc("/leftAt", userHandler.CreateLeftAt).Methods("POST")
	router.HandleFunc("/reserved", reservationHandler.CreateReserved).Methods("POST")
	if len(port) == 0 {
		port = "8080"
	}
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
