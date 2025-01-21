package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"travel-forum-backend/models"
)

// CreateCategory creates a new category
func CreateCategory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Insert the category into the database
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id`
	err = db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// GetCategories retrieves all categories
func GetCategories(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			http.Error(w, "Failed to read categories", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

// GetCategory retrieves a category by ID
func GetCategory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category models.Category
	query := `SELECT id, name, description FROM categories WHERE id = $1`
	err = db.QueryRow(query, categoryID).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

// UpdateCategory updates a category
func UpdateCategory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update the category in the database
	query := `
		UPDATE categories
		SET name = $1, description = $2
		WHERE id = $3`
	_, err = db.Exec(query, category.Name, category.Description, categoryID)
	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory deletes a category
func DeleteCategory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Delete the category from the database
	query := `DELETE FROM categories WHERE id = $1`
	_, err = db.Exec(query, categoryID)
	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}