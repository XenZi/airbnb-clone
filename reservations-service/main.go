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
	router.HandleFunc("/", reservationsHandler.CreateReservation).Methods("POST")
	router.HandleFunc("/availability", reservationsHandler.CreateAvailability).Methods("POST")
	router.HandleFunc("/user/{userId}", reservationsHandler.GetReservationsByUser).Methods("GET")
	router.HandleFunc("/accommodations/reservations/{accommodationID}", reservationsHandler.GetReservationsByAccommodation).Methods("GET")
	router.HandleFunc("/accommodations/reservations", reservationsHandler.ReservationsInDateRangeHandler).Methods("GET")
	router.HandleFunc("/accommodation/dates", reservationsHandler.GetAvailableDates).Methods("GET")
	router.HandleFunc("/{country}/{id}", reservationsHandler.DeleteReservationById).Methods("PUT")

	headersOk := gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := gorillaHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	originsOk := gorillaHandlers.AllowedOrigins([]string{"http://localhost:4200"})
	server := http.Server{
		Addr:         ":" + port,
		Handler:      gorillaHandlers.CORS(headersOk, methodsOk, originsOk)(router),
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
