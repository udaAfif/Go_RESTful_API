package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "Afif Al Amin:Bandarlampung12345@tcp(localhost:3306)/golang test")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	http.HandleFunc("/api/v1/tes", getUsers)
	http.HandleFunc("/api/users", postUser)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

type UserPost struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Location string `json:"location"`
}
type UserGet struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Age      int            `json:"age"`
	Location sql.NullString `json:"location"`
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []UserGet

	for rows.Next() {
		var user UserGet
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.Location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func postUser(w http.ResponseWriter, r *http.Request) {
	var newUser UserPost
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = db.Exec("INSERT INTO users (name, email, age, location) VALUES (?, ?, ?, ?)",
		newUser.Name, newUser.Email, newUser.Age, newUser.Location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT id, name, email, age, location FROM users WHERE email = ?", newUser.Email).Scan(&newUser.ID, &newUser.Name, &newUser.Email, &newUser.Age, &newUser.Location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User added successfully"))
	json.NewEncoder(w).Encode(newUser)

}
