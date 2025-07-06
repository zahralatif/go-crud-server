package main

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        log.Fatal("DB_DSN not found in environment")
    }

	log.Println("DSN from .env:", dsn) //check if DSN is loaded correctly

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal("Failed to open DB:", err)
    }
    defer db.Close()

    if err = db.Ping(); err != nil {
        log.Fatal("DB ping failed:", err)
    }

    log.Println("Connected to PostgreSQL via .env config")
}
