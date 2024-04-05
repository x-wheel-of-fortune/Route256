package main

import (
	"Homework-1/internal/service/service_with_http"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	go func() {
		service_with_http.Secure()
	}()
	service_with_http.Insecure()
}
