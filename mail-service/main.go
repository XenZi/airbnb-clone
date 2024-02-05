package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/XenZi/airbnb-clone/mail-service/config"
	"github.com/XenZi/airbnb-clone/mail-service/domains"
	"github.com/XenZi/airbnb-clone/mail-service/handler"
	"github.com/XenZi/airbnb-clone/mail-service/services"
	"github.com/XenZi/airbnb-clone/mail-service/tracing"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// initialize dependencies
	tracerConfig := tracing.GetConfig()
	tracerProvider, err := tracing.NewTracerProvider("mail-service", tracerConfig.JaegerAddress)
	tracer := tracerProvider.Tracer("mail-service")
	if err != nil {
		log.Fatal(err)
	}
	logger := config.NewLogger("./logs/log.json")
	sender := domains.NewEmailSender(
		os.Getenv("EMAIL"), os.Getenv("PASSWORD"), "smtp.gmail.com", 587)
	mailService := services.NewMailService(sender, logger, tracer)
	mailHandler := handler.NewMailHandler(mailService)

	// router

	router := mux.NewRouter()

	router.HandleFunc("/confirm-new-account", mailHandler.SendAccountConfirmationEmail).Methods("POST")
	router.HandleFunc("/request-reset-password", mailHandler.SendPasswordResetEmail).Methods("POST")
	router.HandleFunc("/send-notification-information", mailHandler.SendNotification).Methods("POST")
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
			logger.LogError("mail-service", err.Error())
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("ERROR PANIC", logrus.Fields{
			"error":  "ERROR PANIC",
			"module": "mail-service main",
		})
	}
	logger.Println("Server stopped")

}
