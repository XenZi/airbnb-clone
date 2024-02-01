package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"recommendation-service/client"
	"recommendation-service/handler"
	"recommendation-service/repository"
	"recommendation-service/services"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sony/gobreaker"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//os env
	port := os.Getenv("PORT")

	// services
	neo4jService, _ := services.NewNeo4j()
	ratingRepository := repository.NewRatingRepository(neo4jService.GetDriver())
	accommodationServiceHost := os.Getenv("ACCOMMODATION_SERVICE_HOST")
	accommodationServicePort := os.Getenv("ACCOMMODATION_SERVICE_PORT")
	userServiceHost := os.Getenv("USER_SERVICE_HOST")
	userServicePort := os.Getenv("USER_SERVICE_PORT")

	customAccommodationClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}
	accommodationServiceCircuitBreaker := gobreaker.NewCircuitBreaker(
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

	customUserClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}
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
	userClient := client.NewUserClient(userServiceHost, userServicePort, customUserClient, userServiceCircuitBreaker)
	accommodationClient := client.NewAccommodationClient(accommodationServiceHost, accommodationServicePort, customAccommodationClient, accommodationServiceCircuitBreaker)
	ratingService := services.NewRatingService(ratingRepository, accommodationClient, userClient)
	ratingHandler := handler.NewRatingHandler(ratingService)
	recommendationRepository := repository.NewRecommendationRepository(neo4jService.GetDriver())
	recommendationService := services.NewRecommendationService(recommendationRepository, accommodationClient)
	recommendationHandler := handler.NewRecommendationHandler(recommendationService)
	// routes

	router := mux.NewRouter()
	router.HandleFunc("/top-rated", recommendationHandler.GetAllRecommendationsByRating).Methods("GET")
	router.HandleFunc("/rating/accommodation/{accommodationID}/{guestID}", ratingHandler.DeleteRatingForAccommodation).Methods("DELETE")
	router.HandleFunc("/rating/host/{id}", ratingHandler.GetAllRatingsForHost).Methods("GET")
	router.HandleFunc("/rating/accommodation/{id}", ratingHandler.GetAllRatingsForAccommmodation).Methods("GET")
	router.HandleFunc("/rating/host", ratingHandler.CreateRatingForHost).Methods("POST")
	router.HandleFunc("/rating/host", ratingHandler.UpdateRatingForHost).Methods("PUT")
	router.HandleFunc("/rating/host/{hostID}/{guestID}", ratingHandler.DeleteRatingForHost).Methods("DELETE")
	router.HandleFunc("/rating/host-by/{hostID}/{guestID}", ratingHandler.GetUserRatingForHost).Methods("GET")
	router.HandleFunc("/rating/accommodation", ratingHandler.CreateRatingForAccommodation).Methods("POST")
	router.HandleFunc("/rating/accommodation", ratingHandler.UpdateRatingForAccommodation).Methods("PUT")
	router.HandleFunc("/{id}", recommendationHandler.GetAllRecommendationsForUser).Methods("GET")
	router.HandleFunc("/rating/{accommodationID}/{guestID}", ratingHandler.GetUserRatingForAccommodation).Methods("GET")
	// server

	if len(port) == 0 {
		port = "8080"
	}
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

	log.Println("Server listening on port", port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	log.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		log.Fatal("Cannot gracefully shutdown...")
	}
	log.Println("Server stopped")
}
