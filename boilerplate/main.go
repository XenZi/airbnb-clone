package main

import (
	"env-test-app/handlers"
	"env-test-app/repository"
	"env-test-app/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	testRepository := repository.TestRepo{}
	testService := services.TestService{Repo: testRepository}
	testHandler := handlers.TestHandler{Service: testService}

	router := mux.NewRouter()

	router.HandleFunc("/test", testHandler.SayHiFromHandler).Methods("GET")

	server := &http.Server{
		Handler: router,
		Addr:    ":8000",
	}
	log.Fatal(server.ListenAndServe())
}
