package main

import (
	"accommodations-service/client"
	"accommodations-service/handlers"
	"accommodations-service/repository"
	"accommodations-service/security"
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
	loggerCach := log.New(os.Stdout, "[cache]", log.LstdFlags)

	//env
	reservationsServiceHost := os.Getenv("RESERVATIONS_SERVICE_HOST")
	log.Println("HOST", reservationsServiceHost)
	reservationsServicePort := os.Getenv("RESERVATIONS_SERVICE_PORT")
	log.Println("PORT", reservationsServicePort)

	userServiceHost := os.Getenv("USER_SERVICE_HOST")
	log.Println("HOST", reservationsServiceHost)
	userServicePort := os.Getenv("USER_SERVICE_PORT")
	log.Println("PORT", reservationsServicePort)

	//clients

	customReservationsServiceClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	customUserServiceClient := &http.Client{
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

	userServiceCircuitBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "user-service",
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

	userClient := client.NewUserClient(userServiceHost, userServicePort, customUserServiceClient, userServiceCircuitBreaker)

	mongoService, err := services.New(timeoutContext, logger)

	if err != nil {
		log.Fatal(err)
	}

	accommodationRepo := repository.NewAccommodationRepository(
		mongoService.GetCli(), logger)
	fileStorage := repository.NewFileStorage(logger)
	defer fileStorage.Close()
	_ = fileStorage.CreateDirectories()
	cache := repository.NewCache(loggerCach)
	accommodationService := services.NewAccommodationService(accommodationRepo, validator, reservationsClient, userClient, fileStorage, cache)
	accommodationsHandler := handlers.AccommodationsHandler{
		AccommodationService: accommodationService,
	}

	accessControl := security.NewAccessControl()
	err = accessControl.LoadAccessConfig("./security/rbac.json")
	if err != nil {
		log.Fatalf("Error loading access configuration: %v", err)
	}
	// ro

	router := mux.NewRouter()

	router.HandleFunc("/recommended", accommodationsHandler.FindAccommodationsByIds).Methods("GET")

	router.HandleFunc("/", accommodationsHandler.GetAllAccommodations).Methods("GET")

	router.HandleFunc("/", accommodationsHandler.CreateAccommodationById).Methods("POST")

	router.HandleFunc("/{id}", accommodationsHandler.UpdateAccommodationById).Methods("PUT")

	router.HandleFunc("/{id}", accommodationsHandler.DeleteAccommodationById).Methods("DELETE")

	router.HandleFunc("/user/{id}", accommodationsHandler.DeleteAccommodationsByUserId).Methods("DELETE")

	router.HandleFunc("/search", accommodationsHandler.SearchAccommodations).Methods("GET")

	router.HandleFunc("/{id}", accommodationsHandler.GetAccommodationById).Methods("GET")

	router.HandleFunc("/images/{id}", accommodationsHandler.MiddlewareCacheHit(accommodationsHandler.GetImage)).Methods("GET")

	router.HandleFunc("/rating/{id}", accommodationsHandler.PutAccommodationRating).Methods("PUT")

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
