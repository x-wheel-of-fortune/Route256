//go:build integration
// +build integration

package tests

import (
	"log"

	"github.com/joho/godotenv"

	"homework/tests/postgresql"
)

var (
	db *postgresql.TDB
)

func init() {
	// тут мы запрашиваем тестовые креды для бд из енв
	// cfg,err := config.FromEnv
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	db = postgresql.NewFromEnv()
}
