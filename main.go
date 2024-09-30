package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

type User struct {
	ID   string `"json":"id"`
	Name string `"json":"name"`
	Age  int    `"json":"age"`
}

var (
	users    []User
	mu       sync.Mutex
	dataFile = "users.json"
)

// Load users from JSON file
func loadUsers() error {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &users)
}

// Save users to JSON file
func saveUsers() error {
	file, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, file, 0644)
}

// Create user
func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter into Create User")

	mu.Lock()
	defer mu.Unlock()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users = append(users, user)
	saveUsers()
	w.WriteHeader(http.StatusCreated)

}

// Update user
func updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter into Update User By ID")

	mu.Lock()
	defer mu.Unlock()

	var updatedUser User
	//json conversion
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, user := range users {
		if user.ID == updatedUser.ID {
			users[i] = updatedUser
			//call save function
			saveUsers()
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

// Get user by ID
func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter into Get User By ID")
	mu.Lock()
	defer mu.Unlock()

	id := r.URL.Query().Get("id")
	for _, user := range users {
		if user.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

// List all users
func listAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter into All Users")

	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func main() {
	// Load users from JSON file
	if err := loadUsers(); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	router := mux.NewRouter()

	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/update", updateUser).Methods("PUT")
	router.HandleFunc("/users/get", getUser).Methods("GET")
	router.HandleFunc("/users/list", listAllUsers).Methods("GET")

	// Start the server
	log.Println(http.ListenAndServe(":8100", router))

}
