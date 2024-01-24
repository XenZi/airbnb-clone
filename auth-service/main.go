package main

import (
	"auth-service/client"
	"auth-service/config"
	"auth-service/handler"
	"auth-service/middlewares"
	"auth-service/repository"
	"auth-service/security"
	"auth-service/services"
	"auth-service/utils"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sony/gobreaker"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger := config.NewLogger("./logs/log.log")

	// env reads

	jwtSecretKey := os.Getenv("JWT_SECRET")
	secretKey := os.Getenv("SECRET_KEY")
	mailServiceHost := os.Getenv("MAIL_SERVICE_HOST")
	mailServicePort := os.Getenv("MAIL_SERVICE_PORT")
	port := os.Getenv("PORT")
	userServiceHost := os.Getenv("USER_SERVICE_HOST")
	userServicePort := os.Getenv("USER_SERVICE_PORT")
	notificationServiceHost := os.Getenv("NOTIFICATION_SERVICE_HOST")
	notificationServicePort := os.Getenv("NOTIFICATION_SERVICE_PORT")
	// clients

	customHttpMailClient := &http.Client{Timeout: time.Second * 10}
	mailClient := client.NewMailClient(mailServiceHost, mailServicePort, customHttpMailClient)

	customUserServiceClient := &http.Client{
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
	userClient := client.NewUserClient(userServiceHost, userServicePort, customUserServiceClient, userServiceCircuitBreaker)

	customNotificationServiceClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	notificationServiceCircuitBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "notification-service",
			MaxRequests: 1,
			Timeout:     10 * time.Second,
			Interval:    0,
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("Circuit Breaker %v: %v -> %v", name, from, to)
			},
		},
	)

	notificationClient := client.NewNotificationClient(notificationServiceHost, notificationServicePort, customNotificationServiceClient, notificationServiceCircuitBreaker)
	// services
	mongoService, err := services.New(timeoutContext, logger)
	if err != nil {
		log.Fatal(err)
	}
	validator := utils.NewValidator()
	userRepo := repository.NewUserRepository(
		mongoService.GetCli(), logger)
	passwordService := services.NewPasswordService()

	keyByte := []byte(jwtSecretKey)
	jwtService := services.NewJWTService(keyByte)
	encryptionService := &services.EncryptionService{SecretKey: secretKey}
	userService := services.NewUserService(userRepo, passwordService, jwtService, validator, encryptionService, mailClient, userClient, notificationClient)
	authHandler := handler.AuthHandler{
		UserService: userService,
	}
	accessControl := security.NewAccessControl()
	err = accessControl.LoadAccessConfig("./security/rbac.json")
	if err != nil {
		log.Fatalf("Error loading access configuration: %v", err)
	}
	// router definitions

	router := mux.NewRouter()
	router.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	router.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/confirm-account/{token}", authHandler.ConfirmAccount).Methods("POST")
	router.HandleFunc("/request-reset-password", authHandler.RequestResetPassword).Methods("POST")
	router.HandleFunc("/reset-password/{token}", authHandler.ResetPassword).Methods("POST")
	router.HandleFunc("/change-password", middlewares.ValidateJWT(authHandler.ChangePassword)).Methods("POST")
	router.HandleFunc("/update-credentials", middlewares.ValidateJWT(authHandler.UpdateCredentials)).Methods("POST")
	router.HandleFunc("/{id}", authHandler.DeleteUser).Methods("DELETE")
	router.HandleFunc("/all", middlewares.ValidateJWT(middlewares.RoleValidator(accessControl, authHandler.All))).Methods("GET")
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
			logger.Panicf("PANIC FROM AUTH-SERVICE ON LISTENING")
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatalf("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")

}
