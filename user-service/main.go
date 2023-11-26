package main

import (
	"context"
	"github.com/gorilla/mux"

	gorillaHandlers "github.com/gorilla/handlers"

	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"user-service/handler"
	"user-service/repository"
	"user-service/service"
	"user-service/utils"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger := log.New(os.Stdout, "[user-api] ", log.LstdFlags)
	mongoService, err := service.New(timeoutContext, logger)
	validator := utils.NewValidator()

	if err != nil {
		log.Fatal(err)
	}
	userRepo := repository.NewUserRepository(mongoService.GetCli(), logger)

	key := os.Getenv("JWT_SECRET")
	keyByte := []byte(key)
	jwtService := service.NewJWTService(keyByte)
	userService := service.NewUserService(userRepo, jwtService, validator)
	profileHandler := handler.UserHandler{
		UserService: userService,
	}
	router := mux.NewRouter()

	router.HandleFunc("/api/create", profileHandler.CreateHandler).Methods("POST")
	router.HandleFunc("/api/update", profileHandler.UpdateHandler).Methods("PUT")
	router.HandleFunc("/api/get-all", profileHandler.GetAllHandler).Methods("GET")
	router.HandleFunc("/api/get-user/{id}", profileHandler.GetUserById).Methods("GET")
	router.HandleFunc("/api/delete/{id}", profileHandler.DeleteHandler).Methods("DELETE")
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

	//Try to shut down gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")

}
