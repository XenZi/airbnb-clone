package main

import (
	"accommodations-service/client"
	"accommodations-service/handlers"
	"accommodations-service/repository"
	"accommodations-service/services"
	"accommodations-service/utils"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sony/gobreaker"

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
	logger := log.New(os.Stdout, "[accommodation-api] ", log.LstdFlags)

	//env
	reservationsServiceHost := os.Getenv("RESERVATIONS_SERVICE_HOST")
	log.Println("HOST", reservationsServiceHost)
	reservationsServicePort := os.Getenv("RESERVATIONS_SERVICE_PORT")
	log.Println("PORT", reservationsServicePort)

	//clients

	customReservationsServiceClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	reservationsServiceCircuitBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "reservations-service",
			MaxRequests: 1,
			Timeout:     10 * time.Second,
			Interval:    0,
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("Circuit Breaker %v: %v -> %v", name, from, to)
			},
		},
	)
	validator := utils.NewValidator()
	reservationsClient := client.NewReservationsClient(reservationsServiceHost, reservationsServicePort, customReservationsServiceClient, reservationsServiceCircuitBreaker)

	mongoService, err := services.New(timeoutContext, logger)

	if err != nil {
		log.Fatal(err)
	}
	accommodationRepo := repository.NewAccommodationRepository(
		mongoService.GetCli(), logger)
	accommodationService := services.NewAccommodationService(accommodationRepo, validator, reservationsClient)
	accommodationsHandler := handlers.AccommodationsHandler{
		AccommodationService: accommodationService,
	}

	router := mux.NewRouter()

	router.HandleFunc("/", accommodationsHandler.GetAllAccommodations).Methods("GET")

	router.HandleFunc("/", accommodationsHandler.CreateAccommodationById).Methods("POST")

	router.HandleFunc("/{id}", accommodationsHandler.UpdateAccommodationById).Methods("PUT")

	router.HandleFunc("/{id}", accommodationsHandler.DeleteAccommodationById).Methods("DELETE")

	router.HandleFunc("/search", accommodationsHandler.SearchAccommodations).Methods("GET")

	router.HandleFunc("/{id}", accommodationsHandler.GetAccommodationById).Methods("GET")

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
