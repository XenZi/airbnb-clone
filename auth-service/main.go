package main

import (
	"auth-service/client"
	"auth-service/handler"
	"auth-service/repository"
	"auth-service/services"
	"auth-service/utils"
	"context"
	"log"
	"net/http"
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

	// env reads

	jwtSecretKey := os.Getenv("JWT_SECRET")
	secretKey := os.Getenv("SECRET_KEY")
	mailServiceHost := os.Getenv("MAIL_SERVICE_HOST")
	mailServicePort := os.Getenv("MAIL_SERVICE_PORT")
	port := os.Getenv("PORT")

	// clients

	customHttpMailClient := &http.Client{Timeout: time.Second * 10}
	mailClient := client.NewMailClient(mailServiceHost, mailServicePort, customHttpMailClient)

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
	userService := services.NewUserService(userRepo, passwordService, jwtService, validator, encryptionService, mailClient)
	authHandler := handler.AuthHandler{
		UserService: userService,
	}

	// router definitions

	router := mux.NewRouter()
	router.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	router.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/confirm-account/{token}", authHandler.ConfirmAccount).Methods("POST")
	router.HandleFunc("/request-reset-password", authHandler.RequestResetPassword).Methods("POST")
	router.HandleFunc("/reset-password/{token}", authHandler.ResetPassword).Methods("POST")

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
