package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/XenZi/airbnb-clone/mail-service/domains"
	"github.com/XenZi/airbnb-clone/mail-service/handler"
	"github.com/XenZi/airbnb-clone/mail-service/services"
	"github.com/gorilla/mux"
)

func main() {

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// initialize dependencies

	logger := log.New(os.Stdout, "[mail-service] ", log.LstdFlags)
	sender := domains.NewEmailSender(
		os.Getenv("EMAIL"), os.Getenv("PASSWORD"), "smtp.gmail.com", 587)
	mailService := services.NewMailService(sender)
	mailHandler := handler.NewMailHandler(mailService)

	// router

	router := mux.NewRouter()

	router.HandleFunc("/confirm-new-account", mailHandler.SendAccountConfirmationEmail).Methods("POST")

	// server

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8081"
	}
	server := http.Server{
		Addr:         ":" + port,
		Handler:      router,
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
