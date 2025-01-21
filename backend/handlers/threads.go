package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"travel-forum-backend/models"
)

// CreateThread creates a new thread
func CreateThread(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var thread models.Thread
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Insert the thread into the database
	query := `
		INSERT INTO threads (title, content, user_id, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	err = db.QueryRow(query, thread.Title, thread.Content, thread.UserID, thread.CategoryID).Scan(&thread.ID, &thread.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create thread", http.StatusInternalServerError)
		print(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(thread)
}

// Get All Threads
func GetAllThreads(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var thread models.Thread
	var err error
	query := `
		SELECT id, title, content, user_id, category_id, created_at
		FROM threads`

	err = db.QueryRow(query).Scan(&thread.ID, &thread.Title, &thread.Content, &thread.UserID, &thread.CategoryID, &thread.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Threads not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(thread)
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

	json.NewEncoder(w).Encode(thread)
}

// UpdateThread updates an existing thread
func UpdateThread(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	var thread models.Thread
	err = json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update the thread in the database
	query := `
		UPDATE threads
		SET title = $1, content = $2, category_id = $3
		WHERE id = $4`
	_, err = db.Exec(query, thread.Title, thread.Content, thread.CategoryID, threadID)
	if err != nil {
		http.Error(w, "Failed to update thread", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(thread)
}

// DeleteThread deletes a thread by ID
func DeleteThread(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	// Delete the thread from the database
	query := `DELETE FROM threads WHERE id = $1`
	_, err = db.Exec(query, threadID)
	if err != nil {
		http.Error(w, "Failed to delete thread", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}