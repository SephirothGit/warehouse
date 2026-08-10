package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	dsn := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("unable to open db connection", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("unable to connect to db", err)
	}

	log.Println("connected to db successfully")
}
