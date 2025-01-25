package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"travel-forum-backend/models"
)

// CreateComment creates a new comment
func CreateComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Check auth_user is the user on the device right now
	authUserID := r.Context().Value("user_id").(int)

	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Set to current user of the device
	comment.UserID = authUserID

	// Insert the comment into the database
	query := `
		INSERT INTO comments (content, user_id, thread_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`
	err = db.QueryRow(query, comment.Content, comment.UserID, comment.ThreadID).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GetCommentsByThread retrieves all comments for a thread
func GetCommentsByThread(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	// Query all comments for the thread
	query := `
		SELECT id, content, user_id, thread_id, created_at
		FROM comments
		WHERE thread_id = $1`
	rows, err := db.Query(query, threadID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.ThreadID, &comment.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to read comments", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// UpdateComment updates an existing comment
func UpdateComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	authUserID := r.Context().Value("user_id").(int)

	var comment models.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Ensure the authenticated user can only update their own comment
	var existingComment models.Comment
	query := `SELECT user_id FROM comments WHERE id = $1`
	err = db.QueryRow(query, commentID).Scan(&existingComment.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Comment not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if existingComment.UserID != authUserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Update the comment in the database
	query = `
		UPDATE comments
		SET content = $1
		WHERE id = $2`
	_, err = db.Exec(query, comment.Content, commentID)
	if err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

// DeleteComment deletes a comment by ID
func DeleteComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	authUserID := r.Context().Value("user_id").(int)

	// Ensure the authenticated user can only delete their own comment
	var existingComment models.Comment
	query := `SELECT user_id FROM comments WHERE id = $1`
	err = db.QueryRow(query, commentID).Scan(&existingComment.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Comment not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if existingComment.UserID != authUserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete the comment from the database
	query = `DELETE FROM comments WHERE id = $1`
	_, err = db.Exec(query, commentID)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}