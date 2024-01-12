package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"recommendation-service/handler"
	"recommendation-service/repository"
	"recommendation-service/services"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//os env
	port := os.Getenv("PORT")

	// services
	neo4jService, _ := services.NewNeo4j()
	ratingRepository := repository.NewRatingRepository(neo4jService.GetDriver())
	ratingService := services.NewRatingService(ratingRepository)
	ratingHandler := handler.NewRatingHandler(ratingService)
	accommodationServiceHost := os.Getenv("ACCOMMODATION_SERVICE_HOST")
	accommodationServicePort := os.Getenv("USER_SERVICE_PORT")
	userServiceHost := os.Getenv("NOTIFICATION_SERVICE_HOST")
	userServicePort := os.Getenv("NOTIFICATION_SERVICE_PORT")

	// routes

	router := mux.NewRouter()
	router.HandleFunc("/rating/host/{id}", ratingHandler.GetAllRatingsForHost).Methods("GET")
	router.HandleFunc("/rating/accommodation/{id}", ratingHandler.GetAllRatingsForAccommmodation).Methods("GET")
	router.HandleFunc("/rating/host", ratingHandler.CreateRatingForHost).Methods("POST")
	router.HandleFunc("/rating/host", ratingHandler.UpdateRatingForHost).Methods("PUT")
	router.HandleFunc("/rating/host", ratingHandler.DeleteRatingForHost).Methods("DELETE")
	router.HandleFunc("/rating/accommodation", ratingHandler.CreateRatingForAccommodation).Methods("POST")
	router.HandleFunc("/rating/accommodation", ratingHandler.UpdateRatingForAccommodation).Methods("PUT")
	router.HandleFunc("/rating/accommodation", ratingHandler.DeleteRatingForAccommodation).Methods("DELETE")

	// server

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
