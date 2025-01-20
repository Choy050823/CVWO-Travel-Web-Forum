package main

import (
	"log"
	"net/http"
	"travel-forum-backend/database"
	"travel-forum-backend/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	database.InitDB()
	defer database.DB.Close()

	// Create a new router
	r := mux.NewRouter()

	// Thread routes
	r.HandleFunc("/threads", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateThread(w, r, database.DB)
	}).Methods("POST")
	r.HandleFunc("/threads/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetThread(w, r, database.DB)
	}).Methods("GET")
	r.HandleFunc("/threads/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateThread(w, r, database.DB)
	}).Methods("PUT")
	r.HandleFunc("/threads/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteThread(w, r, database.DB)
	}).Methods("DELETE")

	// Comment routes
	r.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateComment(w, r, database.DB)
	}).Methods("POST")
	r.HandleFunc("/threads/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCommentsByThread(w, r, database.DB)
	}).Methods("GET")
	r.HandleFunc("/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateComment(w, r, database.DB)
	}).Methods("PUT")
	r.HandleFunc("/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteComment(w, r, database.DB)
	}).Methods("DELETE")

	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}