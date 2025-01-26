package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"travel-forum-backend/models"

	"github.com/gorilla/mux"
)

// CreateThread creates a new thread
func CreateThread(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get user id from jwt
	userID := r.Context().Value("user_id").(int)

	var thread models.Thread
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Set thread user id from JWT
	thread.UserID = userID

	// Validate category_id
	var categoryExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", thread.CategoryID).Scan(&categoryExists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !categoryExists {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Insert the thread into the database
	// Update SQL query to include attached_images
	query := `
		INSERT INTO threads (title, content, user_id, category_id, attached_images)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`
	err = db.QueryRow(query,
		thread.Title,
		thread.Content,
		thread.UserID,
		thread.CategoryID,
		thread.AttachedImages,
	).Scan(&thread.ID, &thread.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create thread", http.StatusInternalServerError)
		print(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(thread)
}

// GetAllThreads retrieves all threads
func GetAllThreads(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query(`
    SELECT id, title, content, user_id, category_id, created_at
    FROM threads
  `)
	if err != nil {
		http.Error(w, "Failed to fetch threads", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Content,
			&thread.UserID,
			&thread.CategoryID, // Ensure this matches the database column name
			&thread.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to read thread data", http.StatusInternalServerError)
			return
		}
		threads = append(threads, thread)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating through threads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(threads)
}

// GetThread retrieves a thread by ID
func GetThread(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	var thread models.Thread
	query := `
		SELECT id, title, content, user_id, category_id, created_at
		FROM threads
		WHERE id = $1`
	err = db.QueryRow(query, threadID).Scan(&thread.ID, &thread.Title, &thread.Content, &thread.UserID, &thread.CategoryID, &thread.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Thread not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(thread)
}

// GetMyThreads retrieves all threads created by the authenticated user
func GetMyThreads(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get user id from JWT
	authUserID := r.Context().Value("user_id").(int)

	// Query threads created by the authenticated user
	query := `
		SELECT id, title, content, user_id, category_id, created_at
		FROM threads
		WHERE user_id = $1`
	rows, err := db.Query(query, authUserID)
	if err != nil {
		http.Error(w, "Failed to fetch threads", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate through the rows and map them to Thread objects
	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		err := rows.Scan(&thread.ID, &thread.Title, &thread.Content, &thread.UserID, &thread.CategoryID, &thread.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to read thread data", http.StatusInternalServerError)
			return
		}
		threads = append(threads, thread)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating through threads", http.StatusInternalServerError)
		return
	}

	// Return the threads as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(threads)
}

func UpdateThread(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	// Get user id from JWT
	authUserID := r.Context().Value("user_id").(int)

	// Fetch the existing thread to check ownership
	var existingThread models.Thread
	query := `SELECT user_id FROM threads WHERE id = $1`
	err = db.QueryRow(query, threadID).Scan(&existingThread.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Thread not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Ensure the authenticated user can only update their own thread
	if existingThread.UserID != authUserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var thread models.Thread
	err = json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate category_id
	var categoryExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", thread.CategoryID).Scan(&categoryExists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !categoryExists {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Update the thread in the database
	query = `
		UPDATE threads
		SET title = $1, content = $2, category_id = $3, attached_images = $4
		WHERE id = $5`
	_, err = db.Exec(query,
		thread.Title,
		thread.Content,
		thread.CategoryID,
		thread.AttachedImages,
		threadID,
	)
	if err != nil {
		http.Error(w, "Failed to update thread", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(thread)
}

func DeleteThread(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	// Get user id from JWT
	authUserID := r.Context().Value("user_id").(int)

	// Fetch the existing thread to check ownership
	var existingThread models.Thread
	query := `SELECT user_id FROM threads WHERE id = $1`
	err = db.QueryRow(query, threadID).Scan(&existingThread.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Thread not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Ensure the authenticated user can only delete their own thread
	if existingThread.UserID != authUserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete the thread from the database
	query = `DELETE FROM threads WHERE id = $1`
	_, err = db.Exec(query, threadID)
	if err != nil {
		http.Error(w, "Failed to delete thread", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func LikeThread(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	query := `UPDATE threads SET likes = likes + 1 WHERE id = $1`
	_, err = db.Exec(query, threadID)
	if err != nil {
		log.Printf("Error liking thread: %v", err)
		http.Error(w, "Failed to like thread", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Thread liked successfully"})
}
