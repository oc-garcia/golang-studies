package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
    Id    sql.NullInt64 `json:"id"`
    Name  string        `json:"name"`
    Email string        `json:"email"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", usersHandler)
	fmt.Println("Server listening on port 3000...")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic(err)
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listUsersHandler(w, r)
	case http.MethodPost:
		createUserHandler(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Name, &u.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert user data
	result, err := db.Exec("INSERT INTO users(name, email) VALUES(?, ?)", u.Name, u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the auto-incremented ID of the newly created user
	var lastInsertId int64
	lastInsertId, err = result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Id = sql.NullInt64{Int64: lastInsertId, Valid: true} // Set Id as sql.NullInt64

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func init() {
	// Connect to the database
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create the "users" table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT
		)
	`)
	if err != nil {
		panic(err)
	}
}