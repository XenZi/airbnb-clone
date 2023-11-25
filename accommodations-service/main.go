package main

import (
	"accommodations-service/handlers"
	"accommodations-service/repository"
	"accommodations-service/services"
	"accommodations-service/utils"
	"context"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	logger := log.New(os.Stdout, "[accommodation-api] ", log.LstdFlags)
	validator := utils.NewValidator()

	mongoService, err := services.New(timeoutContext, logger)

	if err != nil {
		log.Fatal(err)
	}
	accommodationRepo := repository.NewAccommodationRepository(
		mongoService.GetCli(), logger)
	accommodationService := services.NewAccommodationService(accommodationRepo, validator)
	accommodationsHandler := handlers.AccommodationsHandler{
		AccommodationService: accommodationService,
	}

	defer cancel()

	router := mux.NewRouter()

	getAllAccommodations := router.Methods(http.MethodGet).Subrouter()
	getAllAccommodations.HandleFunc("/", accommodationsHandler.GetAllAccommodations)

	getAccommodationsById := router.Methods(http.MethodGet).Subrouter()
	getAccommodationsById.HandleFunc("/{id}", accommodationsHandler.GetAccommodationById)

	postAccommodationForId := router.Methods(http.MethodPost).Subrouter()
	postAccommodationForId.HandleFunc("/", accommodationsHandler.CreateAccommodationById)

	putAccommodationForId := router.Methods(http.MethodPut).Subrouter()
	putAccommodationForId.HandleFunc("/{id}", accommodationsHandler.UpdateAccommodationById)

	deleteAccommodationsById := router.Methods(http.MethodDelete).Subrouter()
	deleteAccommodationsById.HandleFunc("/{id}", accommodationsHandler.DeleteAccommodationById)

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
