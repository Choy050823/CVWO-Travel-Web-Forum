package main

import (
	"log"
	"net/http"
	"travel-forum-backend/database"
	"travel-forum-backend/handlers"
	"travel-forum-backend/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	database.InitDB()
	defer database.DB.Close()

	// Create a new router
	r := mux.NewRouter()


	// Public routes (Any user)
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register(w, r, database.DB)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login(w, r, database.DB)
	}).Methods("POST")

	r.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCategories(w, r, database.DB)
	}).Methods("GET")

	r.HandleFunc("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCategory(w, r, database.DB)
	}).Methods("GET")

	r.HandleFunc("/threads", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllThreads(w, r, database.DB)
	}).Methods("GET")

	r.HandleFunc("/threads/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetThread(w, r, database.DB)
	}).Methods("GET")

	r.HandleFunc("/threads/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCommentsByThread(w, r, database.DB)
	}).Methods("GET")





	// Protected Routes (Login User Only)
	// User routes
	r.HandleFunc("/users/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserProfile(w, r, database.DB)
	})).Methods("GET")

	r.HandleFunc("/users/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUserProfile(w, r, database.DB)
	})).Methods("PUT")

	r.HandleFunc("/users/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteUser(w, r, database.DB)
	})).Methods("DELETE")


	// Categories routes
	r.HandleFunc("/categories", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateCategory(w, r, database.DB)
	})).Methods("POST")

	r.HandleFunc("/categories/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateCategory(w, r, database.DB)
	})).Methods("PUT")

	r.HandleFunc("/categories/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteCategory(w, r, database.DB)
	})).Methods("DELETE")


	// Thread routes
	r.HandleFunc("/threads", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateThread(w, r, database.DB)
	})).Methods("POST")

	r.HandleFunc("/my-threads", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMyThreads(w, r, database.DB)
	}))

	r.HandleFunc("/threads/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateThread(w, r, database.DB)
	})).Methods("PUT")

	r.HandleFunc("/threads/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteThread(w, r, database.DB)
	})).Methods("DELETE")


	// Comment routes
	r.HandleFunc("/comments", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateComment(w, r, database.DB)
	})).Methods("POST")

	r.HandleFunc("/comments/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateComment(w, r, database.DB)
	})).Methods("PUT")

	r.HandleFunc("/comments/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteComment(w, r, database.DB)
	})).Methods("DELETE")





	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
