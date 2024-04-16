package main

import (
	"log"

	"github.com/joho/godotenv"

	"homework/internal/service/service_with_http"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	//go func() {
	//	service_with_http.Secure()
	//}()
	//service_with_http.Insecure()
	service_with_http.Secure()
}
