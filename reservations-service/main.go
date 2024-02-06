package main

import (
	"context"
	"example/saga/messaging/nats"
	"net/http"
	"os"
	"os/signal"
	"reservation-service/client"
	"reservation-service/config"
	"reservation-service/handler"
	"reservation-service/repository"
	"reservation-service/service"
	"reservation-service/utils"
	"time"

	log "github.com/sirupsen/logrus"

	//tracing "command-line-arguments/home/janko33/Documents/airbnb-clone/reservations-service/tracing/tracer.go"

	//opentracing "github.com/opentracing/opentracing-go"

	tracing "reservation-service/tracing"

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
	logger := config.NewLogger("./logs/log.log")
	//storeLogger := log.New(os.Stdout, "[reservation-store]", log.LstdFlags)
	notificationServiceHost := os.Getenv("NOTIFICATION_SERVICE_HOST")
	log.Println("HOST", notificationServiceHost)
	notificationServicePort := os.Getenv("NOTIFICATION_SERVICE_PORT")
	log.Println("PORT", notificationServicePort)
	metricsCommandHost := os.Getenv("COMMAND_SERVICE_HOST")
	log.Println("HOST", metricsCommandHost)
	metricsCommandPort := os.Getenv("COMMAND_SERVICE_PORT")
	log.Println("PORT", metricsCommandPort)
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

	customMetricsServiceClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	metricsServiceCircuitBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "metrics-command",
			MaxRequests: 1,
			Timeout:     10 * time.Second,
			Interval:    0,
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("Circuit Breaker %v: %v -> %v", name, from, to)
			},
		},
	)

	validator := utils.NewValidator()
	notificationsClient := client.NewNotificationClient(notificationServiceHost, notificationServicePort, customNotificationServiceClient, notificationServiceCircuitBreaker)
	metricsClient := client.NewMetricsClient(metricsCommandHost, metricsCommandPort, customMetricsServiceClient, metricsServiceCircuitBreaker)
	tracerConfig := tracing.GetConfig()
	tracerProvider, err := tracing.NewTracerProvider("reservations-service", tracerConfig.JaegerAddress)
	if err != nil {
		log.Fatal("JaegerTraceProvider failed to Initialize", err)
	}
	tracer := tracerProvider.Tracer("reservations-service")

	otel.SetTextMapPropagator(propagation.TraceContext{})
	if err != nil {
		log.Fatal(err)
	}

	store, err := repository.New(logger, tracer)
	if err != nil {
		logger.Fatal("Error while server is listening and serving requests", log.Fields{
			"module": "server-main",
			"error":  err.Error(),
		})
	}
	defer store.CloseSession()
	store.CreateTables()
	reservationRepo, err := repository.New(logger, tracer)
	if err != nil {
		return
	}
	publisher, err := nats.NewNATSPublisher(
		os.Getenv("NATS_HOST"),
		os.Getenv("NATS_PORT"),
		os.Getenv("NATS_USER"),
		os.Getenv("NATS_PASS"),
		os.Getenv("CREATE_ACCOMMODATION_REPLY_SUBJECT"),
	)
	if err != nil {
		log.Fatal(err)
	}
	commandSubscriber, err := nats.NewNATSSubscriber(
		os.Getenv("NATS_HOST"),
		os.Getenv("NATS_PORT"),
		os.Getenv("NATS_USER"),
		os.Getenv("NATS_PASS"),
		os.Getenv("CREATE_ACCOMMODATION_COMMAND_SUBJECT"),
		"reservations-service")
	if err != nil {
		log.Fatal(err)
	}

	reservationService := service.NewReservationService(reservationRepo, validator, notificationsClient, logger, tracer, metricsClient)
	_, err = handler.NewCreateAvailabilityCommandHandler(reservationService, publisher, commandSubscriber, logger)
	if err != nil {
		log.Fatal(err)
	}
	reservationsHandler := handler.ReservationHandler{
		ReservationService: reservationService,
		Tracer:             tracer,
	}
	/*
		tracer, closer := tracing.Init("reservations-service")
		defer closer.Close()
		opentracing.SetGlobalTracer(tracer)


	*/
	router := mux.NewRouter()
	router.HandleFunc("/user/guest/{userId}", reservationsHandler.GetReservationsByUser).Methods("GET")
	router.HandleFunc("/", reservationsHandler.CreateReservation).Methods("POST")
	router.HandleFunc("/accommodations", reservationsHandler.ReservationsInDateRangeHandler).Methods("GET")
	router.HandleFunc("/availability", reservationsHandler.CreateAvailability).Methods("POST")
	router.HandleFunc("/user/host/{hostId}", reservationsHandler.GetReservationsByHost).Methods("GET")
	//router.HandleFunc("/accommodations/{accommodationID}", reservationsHandler.GetReservationsByAccommodation).Methods("GET")
	router.HandleFunc("/accommodation/dates", reservationsHandler.GetAvailableDates).Methods("GET")
	router.HandleFunc("/{country}/{id}/{userID}/{hostID}/{accommodationID}/{endDate}", reservationsHandler.DeleteReservationById).Methods("PUT")
	router.HandleFunc("/{accommodationId}/availability", reservationsHandler.GetAvailabilityForAccommodation).Methods("GET")
	router.HandleFunc("/percentage-cancelation/{hostId}", reservationsHandler.GetCancelationPercentage).Methods("GET")
	router.HandleFunc("/{accommodationId}/{userId}", reservationsHandler.GetReservationsByAccommodationWithEndDate).Methods("GET")
	router.HandleFunc("/host/{hostId}/{userId}", reservationsHandler.GetReservationsByHostWithEndDate).Methods("GET")
	router.HandleFunc("/{accommodationId}/{id}/{country}/{price}", reservationsHandler.UpdateAvailability).Methods("POST")
	router.HandleFunc("/price/myPrice/{maxPrice}", reservationsHandler.GetAccommodationIDsByMaxPrice).Methods("GET")

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
			logger.Fatal("Error while server is listening and serving requests", log.Fields{
				"module": "server-main",
				"error":  err.Error(),
			})
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Error during graceful shutdown", log.Fields{
			"module": "server-main",
			"error":  err.Error(),
		})
	}
	logger.LogInfo("server-main", "Server shut down")

}
