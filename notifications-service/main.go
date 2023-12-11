package main

import (
	"context"
	"log"
	"net/http"
	"notifications-service/handlers"
	"notifications-service/repository"
	"notifications-service/services"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)


func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger := log.New(os.Stdout, "[auth-api] ", log.LstdFlags)

	// env definitions
	port := os.Getenv("PORT")


	// services

	mongoService, err := services.New(timeoutContext, logger)
	if err != nil {
		log.Fatal(err)
	}
	notificationRepository := repository.NewNotificationRepository(mongoService.GetCli(), logger)
	notificationService := services.NewNotificationService(notificationRepository)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	
	
	// router definitions

	router := mux.NewRouter()
	router.HandleFunc("/create-new-user-notification/{id}", notificationHandler.CreateNewUserNotification).Methods("POST")
	router.HandleFunc("/{id}", notificationHandler.CreateNewNotificationForUser).Methods("POST")
	router.HandleFunc("/{id}", notificationHandler.ReadAllNotifications).Methods("PUT")

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