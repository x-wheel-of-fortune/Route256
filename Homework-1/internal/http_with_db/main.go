package main

import (
	"Homework-1/internal/service"
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
		service.Secure()
	}()
	service.Insecure()
}
