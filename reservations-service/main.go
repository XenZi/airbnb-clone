package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reservation-service/client"
	"reservation-service/handler"
	"reservation-service/repository"
	"reservation-service/service"
	"reservation-service/tracing"
	"reservation-service/utils"
	"time"

	//tracing "command-line-arguments/home/janko33/Documents/airbnb-clone/reservations-service/tracing/tracer.go"

	//opentracing "github.com/opentracing/opentracing-go"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	notificationServiceHost := os.Getenv("NOTIFICATION_SERVICE_HOST")
	log.Println("HOST", notificationServiceHost)
	notificationServicePort := os.Getenv("NOTIFICATION_SERVICE_PORT")
	log.Println("PORT", notificationServicePort)

	customNotificationServiceClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	notificationServiceCircuitBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "notifications-service",
			MaxRequests: 1,
			Timeout:     10 * time.Second,
			Interval:    0,
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("Circuit Breaker %v: %v -> %v", name, from, to)
			},
		},
	)
	cfg := tracing.GetConfig()

	tracerProvider, err := tracing.NewTracerProvider(cfg.ServiceName, cfg.JaegerAddress)
	if err != nil {
		log.Fatal("JaegerTraceProvider failed to Initialize", err)
	}
	tracer := tracerProvider.Tracer(cfg.ServiceName)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	validator := utils.NewValidator()
	notificationsClient := client.NewNotificationClient(notificationServiceHost, notificationServicePort, customNotificationServiceClient, notificationServiceCircuitBreaker, tracer)

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

	/*	ctx := context.Background()
		exp, err := tracing.NewExporter(cfg.JaegerAddress)
		if err != nil {
			log.Fatalf("failed to initialize exporter: %v", err)
		}*/
	// Create a new tracer provider with a batch span processor and the given exporter.
	//	tp := tracing.NewTraceProvider(exp)
	// Handle shutdown properly so nothing leaks.
	//	defer func() { _ = tp.Shutdown(ctx) }()
	//	otel.SetTracerProvider(tp)
	// Finally, set the tracer that can be used for this package.
	//	tracer := tp.Tracer("ordering-service")

	otel.SetTextMapPropagator(propagation.TraceContext{})
	reservationService := service.NewReservationService(reservationRepo, validator, notificationsClient, tracer)
	reservationsHandler := handler.ReservationHandler{
		Logger:             logger,
		Tracer:             tracer,
		ReservationService: reservationService,
	}
	/*
		tracer, closer := tracing.Init("reservations-service")
		defer closer.Close()
		opentracing.SetGlobalTracer(tracer)


	*/
	router := mux.NewRouter()
	router.Use(handler.ExtractTraceInfoMiddleware)
	router.HandleFunc("/user/guest/{userId}", reservationsHandler.GetReservationsByUser).Methods("GET")
	router.HandleFunc("/", reservationsHandler.CreateReservation).Methods("POST")
	router.HandleFunc("/accommodations", reservationsHandler.ReservationsInDateRangeHandler).Methods("GET")
	router.HandleFunc("/availability", reservationsHandler.CreateAvailability).Methods("POST")
	router.HandleFunc("/user/host/{hostId}", reservationsHandler.GetReservationsByHost).Methods("GET")
	//router.HandleFunc("/accommodations/{accommodationID}", reservationsHandler.GetReservationsByAccommodation).Methods("GET")
	router.HandleFunc("/accommodation/dates", reservationsHandler.GetAvailableDates).Methods("GET")
	router.HandleFunc("/{country}/{id}/{userId}/{hostId}/{accommodationId}", reservationsHandler.DeleteReservationById).Methods("DELETE")
	router.HandleFunc("/{accommodationId}/availability", reservationsHandler.GetAvailabilityForAccommodation).Methods("GET")
	router.HandleFunc("/percentage-cancelation/{hostId}", reservationsHandler.GetCancelationPercentage).Methods("GET")
	router.HandleFunc("/{accommodationId}/{userId}", reservationsHandler.GetReservationsByAccommodationWithEndDate).Methods("GET")
	router.HandleFunc("/host/{hostId}/{userId}", reservationsHandler.GetReservationsByHostWithEndDate).Methods("GET")

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
