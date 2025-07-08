package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"

    _ "github.com/lib/pq"
)

type Employee struct {
    ID         int     `json:"id"`
    Name       string  `json:"name"`
    Email      string  `json:"email"`
    Position   string  `json:"position"`
    Department string  `json:"department"`
    Salary     float64 `json:"salary"`
    CreatedAt  string  `json:"created_at"`
}

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

http.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        // GET /employees - all employees
        rows, err := db.Query("SELECT id, name, email, position, department, salary, created_at FROM employees")
        if err != nil {
            http.Error(w, "DB query failed", http.StatusInternalServerError)
            log.Println("Query error:", err)
            return
        }
        defer rows.Close()

        var employees []Employee
        for rows.Next() {
            var emp Employee
            err := rows.Scan(&emp.ID, &emp.Name, &emp.Email, &emp.Position, &emp.Department, &emp.Salary, &emp.CreatedAt)
            if err != nil {
                http.Error(w, "Scan failed", http.StatusInternalServerError)
                log.Println("Scan error:", err)
                return
            }
            employees = append(employees, emp)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(employees)

    case http.MethodPost:
        // POST /employees - create new employee
        var emp Employee
        err := json.NewDecoder(r.Body).Decode(&emp)
        if err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            log.Println("Decode error:", err)
            return
        }

        query := `INSERT INTO employees (name, email, position, department, salary)
                  VALUES ($1, $2, $3, $4, $5) RETURNING id`
        err = db.QueryRow(query, emp.Name, emp.Email, emp.Position, emp.Department, emp.Salary).Scan(&emp.ID)
        if err != nil {
            http.Error(w, "DB insert failed", http.StatusInternalServerError)
            log.Println("Insert error:", err)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(emp)

    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
})

// Handle individual employee operations
http.HandleFunc("/employees/", func(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/employees/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    switch r.Method {
    case http.MethodGet:
        // GET /employees/{id}
        var emp Employee
        query := `SELECT id, name, email, position, department, salary, created_at FROM employees WHERE id = $1`
        row := db.QueryRow(query, id)
        err := row.Scan(&emp.ID, &emp.Name, &emp.Email, &emp.Position, &emp.Department, &emp.Salary, &emp.CreatedAt)
        if err == sql.ErrNoRows {
            http.Error(w, "Employee not found", http.StatusNotFound)
            return
        } else if err != nil {
            http.Error(w, "DB error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(emp)

    case http.MethodPut:
        // PUT /employees/{id}
        var emp Employee
        err := json.NewDecoder(r.Body).Decode(&emp)
        if err != nil {
            http.Error(w, "Invalid JSON body", http.StatusBadRequest)
            return
        }

        query := `UPDATE employees SET name=$1, email=$2, position=$3, department=$4, salary=$5 WHERE id=$6`
        res, err := db.Exec(query, emp.Name, emp.Email, emp.Position, emp.Department, emp.Salary, id)
        if err != nil {
            http.Error(w, "DB update failed", http.StatusInternalServerError)
            return
        }

        count, err := res.RowsAffected()
        if err != nil || count == 0 {
            http.Error(w, "No employee updated", http.StatusNotFound)
            return
        }

        emp.ID = id
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(emp)

    // DELETE /employees/{id}
    case http.MethodDelete:
        query := `DELETE FROM employees WHERE id=$1`
        res, err := db.Exec(query, id)
        if err != nil {
            http.Error(w, "DB delete failed", http.StatusInternalServerError)
            return
        }

        count, err := res.RowsAffected()
        if err != nil || count == 0 {
            http.Error(w, "Employee not found", http.StatusNotFound)
            return
        }

        w.WriteHeader(http.StatusNoContent)

    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
})

    // Start the HTTP server
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("HTTP server failed: %v", err)
    }
}
