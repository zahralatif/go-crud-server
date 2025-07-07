package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    _ "github.com/lib/pq"
)

const dbConnStr = "postgres://zahra:Zl-l%40b-l%40t1%21%40%23@db-01.lab.internal/amhsdb?sslmode=disable"

func main() {
    db, err := sql.Open("postgres", dbConnStr)
    if err != nil {
        log.Fatalf("DB connection failed: %v", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatalf("DB ping failed: %v", err)
    }

    fmt.Println("Connected to DB")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Go server running!")
    })

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("HTTP server failed: %v", err)
    }
}

// File: main.go
// Description: This Go program connects to a PostgreSQL database using a DSN loaded from a
// .env file. It uses the "github.com/joho/godotenv" package to load environment variables
// and the "github.com/lib/pq" package as the PostgreSQL driver.
// It logs the connection status and handles errors appropriately.
// Usage: Ensure you have a .env file with the DB_DSN variable set to your PostgreSQL connection string.


// package main

// import (
//     "database/sql"
//     "log"
//     "os"

//     _ "github.com/lib/pq"
//     "github.com/joho/godotenv"
// )

// func main() {
//     err := godotenv.Load()
//     if err != nil {
//         log.Fatal("Error loading .env file")
//     }

//     dsn := os.Getenv("DB_DSN")
//     if dsn == "" {
//         log.Fatal("DB_DSN not found in environment")
//     }

// 	log.Println("DSN from .env:", dsn) //check if DSN is loaded correctly

//     db, err := sql.Open("postgres", dsn)
//     if err != nil {
//         log.Fatal("Failed to open DB:", err)
//     }
//     defer db.Close()

//     if err = db.Ping(); err != nil {
//         log.Fatal("DB ping failed:a", err)
//     }

//     log.Println("Connected to PostgreSQL via .env config")
// }
