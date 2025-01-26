package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"travel-forum-backend/models"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

// CreateComment creates a new comment
// CreateComment creates a new comment
func CreateComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Use the Comment model from your models package
	var comment models.Comment

	// Decode the request body into the Comment model
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if comment.Content == "" || comment.UserID == 0 || comment.ThreadID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Fetch the username from the users table
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id = $1", comment.UserID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Error fetching username: %v", err)
			http.Error(w, "Failed to fetch username", http.StatusInternalServerError)
		}
		return
	}

	// Insert the comment into the database
	query := `
        INSERT INTO comments (content, user_id, thread_id, attached_images, upvotes, downvotes)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at`
	err = db.QueryRow(query,
		comment.Content,
		comment.UserID,
		comment.ThreadID,
		pq.Array(comment.AttachedImages), // Use pq.Array for PostgreSQL arrays
		comment.Upvotes,
		comment.Downvotes,
	).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		log.Printf("Error creating comment: %v", err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	// Add the username to the response
	comment.Author = username

	// Return the created comment
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

	query := `
        SELECT c.id, c.content, c.user_id, c.thread_id, c.created_at, 
               c.attached_images, c.upvotes, c.downvotes, u.username
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.thread_id = $1
        ORDER BY c.created_at DESC`

	// Log the query and thread ID for debugging
	log.Printf("Executing query for thread ID: %d", threadID)

	rows, err := db.Query(query, threadID)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var username string
		var attachedImages []string

		// Scan the row into variables
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.UserID,
			&comment.ThreadID,
			&comment.CreatedAt,
			pq.Array(&attachedImages), // Use pq.Array for PostgreSQL arrays
			&comment.Upvotes,
			&comment.Downvotes,
			&username,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Failed to read comments", http.StatusInternalServerError)
			return
		}

		// Assign the scanned values to the comment struct
		comment.AttachedImages = attachedImages
		comment.Author = username
		comments = append(comments, comment)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		http.Error(w, "Failed to read comments", http.StatusInternalServerError)
		return
	}

	// Return the comments as JSON
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

// Add VoteComment handler
// func VoteComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	vars := mux.Vars(r)
// 	commentID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
// 		return
// 	}

// 	var vote struct {
// 		Type string `json:"type"` // "up" or "down"
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	var query string
// 	switch vote.Type {
// 	case "up":
// 		query = "UPDATE comments SET upvotes = upvotes + 1 WHERE id = $1"
// 	case "down":
// 		query = "UPDATE comments SET downvotes = downvotes + 1 WHERE id = $1"
// 	default:
// 		http.Error(w, "Invalid vote type", http.StatusBadRequest)
// 		return
// 	}

// 	_, err = db.Exec(query, commentID)
// 	if err != nil {
// 		http.Error(w, "Failed to update vote", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// handlers/comments.go
func UpvoteComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    vars := mux.Vars(r)
    commentID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid comment ID", http.StatusBadRequest)
        return
    }

    query := `UPDATE comments SET upvotes = upvotes + 1 WHERE id = $1`
    _, err = db.Exec(query, commentID)
    if err != nil {
        log.Printf("Error upvoting comment: %v", err)
        http.Error(w, "Failed to upvote comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Comment upvoted successfully"})
}

func DownvoteComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    vars := mux.Vars(r)
    commentID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid comment ID", http.StatusBadRequest)
        return
    }

    query := `UPDATE comments SET downvotes = downvotes + 1 WHERE id = $1`
    _, err = db.Exec(query, commentID)
    if err != nil {
        log.Printf("Error downvoting comment: %v", err)
        http.Error(w, "Failed to downvote comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Comment downvoted successfully"})
}