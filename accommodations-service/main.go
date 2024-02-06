package main

import (
	"accommodations-service/client"
	"accommodations-service/config"
	"accommodations-service/handlers"
	"accommodations-service/middlewares"
	"accommodations-service/orchestrator"
	"accommodations-service/repository"
	"accommodations-service/security"
	"accommodations-service/services"
	"accommodations-service/tracing"
	"accommodations-service/utils"
	"context"
	"example/saga/messaging/nats"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

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
	loggerW := config.NewLogger("./logs/log.log")

	loggerCach := log.New(os.Stdout, "[cache]", log.LstdFlags)

	//env
	reservationsServiceHost := os.Getenv("RESERVATIONS_SERVICE_HOST")
	log.Println("HOST", reservationsServiceHost)
	reservationsServicePort := os.Getenv("RESERVATIONS_SERVICE_PORT")
	log.Println("PORT", reservationsServicePort)

	userServiceHost := os.Getenv("USER_SERVICE_HOST")
	log.Println("HOST", userServiceHost)
	userServicePort := os.Getenv("USER_SERVICE_PORT")
	log.Println("PORT", userServicePort)

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
	reservationsClient := client.NewReservationsClient(reservationsServiceHost, reservationsServicePort, customReservationsServiceClient, reservationsServiceCircuitBreaker, loggerW)
	userClient := client.NewUserClient(userServiceHost, userServicePort, customUserServiceClient, userServiceCircuitBreaker, loggerW)

	tracerConfig := tracing.GetConfig()
	tracerProvider, err := tracing.NewTracerProvider("accommodations-service", tracerConfig.JaegerAddress)
	if err != nil {
		log.Fatal("JaegerTraceProvider failed to Initialize", err)
	}
	tracer := tracerProvider.Tracer("accommodations-service")

	otel.SetTextMapPropagator(propagation.TraceContext{})
	if err != nil {
		log.Fatal(err)
	}

	mongoService, err := services.New(timeoutContext, loggerW)

	if err != nil {
		log.Fatal(err)
	}

	accommodationRepo := repository.NewAccommodationRepository(
		mongoService.GetCli(), loggerW, tracer)
	publisher, err := nats.NewNATSPublisher(
		os.Getenv("NATS_HOST"),
		os.Getenv("NATS_PORT"),
		os.Getenv("NATS_USER"),
		os.Getenv("NATS_PASS"),
		os.Getenv("CREATE_ACCOMMODATION_COMMAND_SUBJECT"),
	)
	if err != nil {
		log.Fatal(err)
	}
	replySubscriber, err := nats.NewNATSSubscriber(
		os.Getenv("NATS_HOST"),
		os.Getenv("NATS_PORT"),
		os.Getenv("NATS_USER"),
		os.Getenv("NATS_PASS"),
		os.Getenv("CREATE_ACCOMMODATION_REPLY_SUBJECT"),
		"accommodations-service")
	if err != nil {
		log.Fatal(err)
	}
	orch, err := orchestrator.NewCreateAccommodationOrchestrator(publisher, replySubscriber, loggerW)
	if err != nil {
		log.Fatal(err)
	}
	fileStorage := repository.NewFileStorage(loggerW, tracer)
	defer fileStorage.Close()
	_ = fileStorage.CreateDirectories()
	cache := repository.NewCache(loggerCach, tracer)
	accommodationService := services.NewAccommodationService(accommodationRepo, validator, reservationsClient, userClient, fileStorage, cache, orch, tracer, loggerW)
	publisher1, err := nats.NewNATSPublisher(
		os.Getenv("NATS_HOST"),
		os.Getenv("NATS_PORT"),
		os.Getenv("NATS_USER"),
		os.Getenv("NATS_PASS"),
		os.Getenv("CREATE_ACCOMMODATION_REPLY_SUBJECT"),
	)
	if err != nil {
		log.Fatal(err)
	}
	replySubscriber2, err := nats.NewNATSSubscriber(
		os.Getenv("NATS_HOST"),
		os.Getenv("NATS_PORT"),
		os.Getenv("NATS_USER"),
		os.Getenv("NATS_PASS"),
		os.Getenv("CREATE_ACCOMMODATION_COMMAND_SUBJECT"),
		"accommodations-service")
	_, err = handlers.NewCreateAccommodationCommandHandler(accommodationService, publisher1, replySubscriber2, tracer, loggerW)
	if err != nil {
		log.Println(err)
	}

	accommodationsHandler := handlers.AccommodationsHandler{
		AccommodationService: accommodationService,
		Tracer:               tracer,
		Logger:               loggerW,
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

	router.HandleFunc("/", middlewares.ValidateJWT(middlewares.RoleValidator("Host", accommodationsHandler.CreateAccommodationById))).Methods("POST")

	router.HandleFunc("/{id}", middlewares.ValidateJWT(middlewares.RoleValidator("Host", accommodationsHandler.UpdateAccommodationById))).Methods("PUT")

	router.HandleFunc("/{id}", middlewares.ValidateJWT(middlewares.RoleValidator("Host", accommodationsHandler.DeleteAccommodationById))).Methods("DELETE")

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

	loggerW.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			loggerW.Fatalf(err.Error())
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	loggerW.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		loggerW.Fatalf("Cannot gracefully shutdown...")
	}
	loggerW.Println("Server stopped")

}
