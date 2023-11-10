package main

import (
	"auth-service/handler"
	"auth-service/repository"
	"auth-service/services"
	"auth-service/utils"
	"context"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger := log.New(os.Stdout, "[auth-api] ", log.LstdFlags)
	mongoService, err := services.New(timeoutContext, logger)
	validator := utils.NewValidator()

	if err != nil {
		log.Fatal(err)
	}
	userRepo := repository.NewUserRepository(
		mongoService.GetCli(), logger)
	passwordService := services.NewPasswordService()

	key := os.Getenv("JWT_SECRET")
	keyByte := []byte(key)
	jwtService := services.NewJWTService(keyByte)
	userService := services.NewUserService(userRepo, passwordService, jwtService, validator)
	authHandler := handler.AuthHandler{
		UserService: userService,
	}
	router := mux.NewRouter()

	router.HandleFunc("/api/login", authHandler.LoginHandler).Methods("POST")
	router.HandleFunc("/api/register", authHandler.RegisterHandler).Methods("POST")
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
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
