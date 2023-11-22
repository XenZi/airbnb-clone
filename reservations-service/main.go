package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reservation-service/handler"
	"reservation-service/repository"
	"reservation-service/service"
	"reservation-service/utils"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[reservation-api]", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[reservation-store]", log.LstdFlags)
	validator := utils.NewValidator()

	store, err := repository.New(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.CloseSession()
	store.CreateTables()
	reservationRepo, err := repository.New(logger)
	if err != nil {
		return
	}
	reservationService := service.NewReservationService(reservationRepo, validator)
	reservationsHandler := handler.ReservationHandler{
		ReservationService: reservationService,
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/reservations", reservationsHandler.CreateReservationByUser).Methods("POST")
	router.HandleFunc("/api/reservations/user/{userId}", reservationsHandler.GetReservationsByUser).Methods("GET")
	router.HandleFunc("/api/reservations/{userId}/{id}", reservationsHandler.DeleteReservationById).Methods("DELETE")

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")

}
