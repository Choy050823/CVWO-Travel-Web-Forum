package main

import (
	"log"
	"net/http"
	"travel-forum-backend/database"
	"travel-forum-backend/handlers"
	"travel-forum-backend/middleware"
	"github.com/rs/cors"
	"github.com/gorilla/mux"
)

// func enableCORS(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174") // Replace with your React app's URL
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		if r.Method == "OPTIONS" {
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

func main() {
	// Initialize the database
	database.InitDB()
	defer database.DB.Close()

	// Create a new router
	r := mux.NewRouter()
	// r.Use(enableCORS)

	// Public routes (Any user)
	r.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register(w, r, database.DB)
	}).Methods("POST")

	r.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login(w, r, database.DB)
	}).Methods("POST")

	r.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCategories(w, r, database.DB)
	}).Methods("GET")

	r.HandleFunc("/api/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCategory(w, r, database.DB)
	}).Methods("GET")

	r.HandleFunc("/api/threads", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllThreads(w, r, database.DB)
	}).Methods("GET")

	r.HandleFunc("/api/threads/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetThread(w, r, database.DB)
	}).Methods("GET")

	r.HandleFunc("/api/threads/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCommentsByThread(w, r, database.DB)
	}).Methods("GET")

	// Protected Routes (Login User Only)
	// User routes
	r.HandleFunc("/api/users/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserProfile(w, r, database.DB)
	})).Methods("GET")

	r.HandleFunc("/api/users/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUserProfile(w, r, database.DB)
	})).Methods("PUT")

	r.HandleFunc("/api/users/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteUser(w, r, database.DB)
	})).Methods("DELETE")

	// Categories routes
	r.HandleFunc("/api/categories", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateCategory(w, r, database.DB)
	})).Methods("POST")

	r.HandleFunc("/api/categories/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateCategory(w, r, database.DB)
	})).Methods("PUT")

	r.HandleFunc("/api/categories/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteCategory(w, r, database.DB)
	})).Methods("DELETE")

	// Thread routes
	r.HandleFunc("/api/threads", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateThread(w, r, database.DB)
	})).Methods("POST")

	r.HandleFunc("/api/my-threads", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMyThreads(w, r, database.DB)
	}))

	r.HandleFunc("/api/threads/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateThread(w, r, database.DB)
	})).Methods("PUT")

	r.HandleFunc("/api/threads/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteThread(w, r, database.DB)
	})).Methods("DELETE")

	// Comment routes
	r.HandleFunc("/api/comments", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateComment(w, r, database.DB)
	})).Methods("POST")

	r.HandleFunc("/api/comments/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateComment(w, r, database.DB)
	})).Methods("PUT")

	r.HandleFunc("/api/comments/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteComment(w, r, database.DB)
	})).Methods("DELETE")

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5174"}, // Allow your React frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap your router with the CORS middleware
	handler := c.Handler(r)

	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
