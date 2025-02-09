package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"travel-forum-backend/cache"
	"travel-forum-backend/models"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
// var jwtKey = []byte("5lfX8Bl4C1mZZ/ljU+BrWFoxTcxQqacwPVfloDs+5No=")

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate email
	if !isValidEmail(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validate password strength
	if len(req.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	var existingUser models.User
	err = db.QueryRow(`SELECT id FROM users WHERE email = $1`, req.Email).Scan(&existingUser.ID)
	if err == nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		log.Printf("Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check if username already exists
	err = db.QueryRow(`SELECT id FROM users WHERE username = $1`, req.Username).Scan(&existingUser.ID)
	if err == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		log.Printf("Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
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
		log.Printf("Database error: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Return the user (without password hash)
	user.Username = req.Username
	user.Email = req.Email
	user.Role = "user"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

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
    fmt.Println("Executing query:", query) // Log the query

    err = db.QueryRow(query, req.Email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
    if err != nil {
        fmt.Println("Query error:", err) // Log the error from QueryRow
        if err == sql.ErrNoRows {
            http.Error(w, "User not found", http.StatusNotFound)
        } else {
            http.Error(w, "Database error", http.StatusInternalServerError)
        }
        return
    }

    fmt.Println("User found:", user) // Log the user data (for debugging only - remove in production)

	// Compare the password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create JWT claims
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Generate the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Cache the user credential
	cache.CacheUser(strconv.Itoa(user.ID), user)

	// Return the token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract the claims from the request context
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		log.Println("Invalid token claims in GetCurrentUser")
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	log.Printf("Fetching user details for UserID: %d", claims.UserID)

	// Fetch the user from the database
	var user models.User
	query := `SELECT id, username, email, role FROM users WHERE id = $1`
	err := db.QueryRow(query, claims.UserID).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not found for UserID: %d", claims.UserID)
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Database error in GetCurrentUser: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("User details fetched successfully: %+v", user)

	// Return the user data
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
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

	// Check auth_user is the user on the device right now
	authUserID := r.Context().Value("user_id").(int)

	if authUserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	updatedUser := models.User{
		ID:       authUserID,
		Username: req.Username,
		Email:    req.Email,
	}
	// Cache the new user
	cache.CacheUser(strconv.Itoa(authUserID), updatedUser)

	w.Header().Set("Content-Type", "application/json")
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

	// Check auth_user is the user on the device right now
	authUserID := r.Context().Value("user_id").(int)

	if authUserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete the user from the database
	query := `DELETE FROM users WHERE id = $1`
	_, err = db.Exec(query, userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Remove the user from the cache
	cache.DeleteCachedUser((string)(userID))

	w.WriteHeader(http.StatusNoContent)
}
