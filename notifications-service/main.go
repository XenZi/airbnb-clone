package main

import (
	"context"
	"log"
	"net/http"
	"notifications-service/client"
	"notifications-service/handlers"
	"notifications-service/repository"
	"notifications-service/services"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sony/gobreaker"
)


func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger := log.New(os.Stdout, "[auth-api] ", log.LstdFlags)

	// env definitions
	port := os.Getenv("PORT")
	mailHost := os.Getenv("MAIL_SERVICE_HOST")
	mailPort := os.Getenv("MAIL_SERVICE_PORT") 
	userServiceHost := os.Getenv("USER_SERVICE_HOST")
	userServicePort := os.Getenv("USER_SERVICE_PORT")

	// commms

	customMailClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost: 10,
		},
	}
	mailCircuitBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name: "mail-service",
			MaxRequests: 1,
			Timeout: 10 * time.Second,
			Interval: 0,
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("Circuit Breaker %v: %v -> %v", name, from, to)
			},
		},
	)
	mailClient := client.NewMailClient(mailHost, mailPort, customMailClient, mailCircuitBreaker)

	customUserServiceClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost: 10,
		},
	}
	
	userServiceCircuitBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name: "user-service",
			MaxRequests: 1,
			Timeout: 10 * time.Second,
			Interval: 0,
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("Circuit Breaker %v: %v -> %v", name, from, to)
			},
		},
	)
	userClient := client.NewUserClient(userServiceHost, userServicePort, customUserServiceClient, userServiceCircuitBreaker)

	// services

	mongoService, err := services.New(timeoutContext, logger)
	if err != nil {
		log.Fatal(err)
	}
	notificationRepository := repository.NewNotificationRepository(mongoService.GetCli(), logger)
	notificationService := services.NewNotificationService(notificationRepository, mailClient, userClient)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	
	
	// router definitions

	router := mux.NewRouter()
	router.HandleFunc("/create-new-user-notification/{id}", notificationHandler.CreateNewUserNotification).Methods("POST")
	router.HandleFunc("/{id}", notificationHandler.CreateNewNotificationForUser).Methods("POST")
	router.HandleFunc("/{id}", notificationHandler.ReadAllNotifications).Methods("PUT")
	router.HandleFunc("/{id}", notificationHandler.GetAllNotificationsByID).Methods("GET")

	// server definitions

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

	logger.Println("Server listening on port", port)

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