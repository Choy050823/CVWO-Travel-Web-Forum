package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	// "time"

	"travel-forum-backend/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert the user into the database
	query := `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	var user models.User
	err = db.QueryRow(query, req.Username, req.Email, hashedPassword, "user").Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Return the user (without password hash)
	user.Username = req.Username
	user.Email = req.Email
	user.Role = "user"

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Fetch the user from the database
	var user models.User
	query := `SELECT id, username, email, password_hash, role, created_at FROM users WHERE email = $1`
	err = db.QueryRow(query, req.Email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Compare the password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Return the user (without password hash)
	user.PasswordHash = ""
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetUserProfile retrieves a user's profile
func GetUserProfile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch the user from the database
	var user models.User
	query := `SELECT id, username, email, role, created_at FROM users WHERE id = $1`
	err = db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateUserProfile updates a user's profile
func UpdateUserProfile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req models.User
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update the user in the database
	query := `
		UPDATE users
		SET username = $1, email = $2
		WHERE id = $3`
	_, err = db.Exec(query, req.Username, req.Email, userID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

// DeleteUser deletes a user
func DeleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Delete the user from the database
	query := `DELETE FROM users WHERE id = $1`
	_, err = db.Exec(query, userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}