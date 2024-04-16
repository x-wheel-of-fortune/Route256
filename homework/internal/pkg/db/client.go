package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDb(ctx context.Context) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, generateDsn())
	if err != nil {
		return nil, err
	}
	return newDatabase(pool), nil
}

func generateDsn() string {
	host, exists := os.LookupEnv("HOST")
	if !exists {
		log.Println("Не указан HOST")
	}
	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Println("Не указан PORT")
	}
	user, exists := os.LookupEnv("USER")
	if !exists {
		log.Println("Не указан USER")
	}
	password, exists := os.LookupEnv("PASSWORD")
	if !exists {
		log.Println("Не указан PASSWORD")
	}
	dbname, exists := os.LookupEnv("DBNAME")
	if !exists {
		log.Println("Не указано DBNAME")
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}
