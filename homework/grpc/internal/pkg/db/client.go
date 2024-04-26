package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDb(ctx context.Context) (*Database, error) {
	var err error
	for i := 0; i < 5; i++ {
		s := generateDsn()
		pool, err := pgxpool.Connect(ctx, s)
		if err != nil {
			time.Sleep(2 * time.Second)
			log.Printf("Failed to connect to database %s, %d/5\n", s, i+1)
		} else {
			return newDatabase(pool), nil
		}
	}
	return nil, err
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
