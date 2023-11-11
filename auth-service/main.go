package main

import (
	"auth-service/handler"
	"auth-service/repository"
	"auth-service/services"
	"auth-service/utils"
	"context"
	"github.com/XenZi/airbnb-clone/api-gateway/proto/auth_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	listener, err := net.Listen("tcp", os.Getenv("AUTH_SERVICE_ADDRESS"))

	if err != nil {
		log.Fatalln(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

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
	auth_service.RegisterAuthServiceServer(grpcServer, authHandler)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	grpcServer.Stop()

}
