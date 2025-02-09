package main

import (
	"context"

	// "path/filepath"
	// "database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"travel-forum-backend/database"
	"travel-forum-backend/handlers"
	"travel-forum-backend/middleware"

	"github.com/gorilla/mux"
	// "github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/rs/cors"

	// "github.com/golang-migrate/migrate/v4" // Add this import

    _ "github.com/golang-migrate/migrate/v4/database/postgres" // Import the postgres driver
)

func main() {
	// local test only
	// Construct path to .env file in the root directory
	// rootDir, err := os.Getwd()
	// if err != nil {
	// 		log.Fatal("Error getting working directory:", err)
	// }

	// envPath := filepath.Join(rootDir, "..", ".env") // Go up one level for .env
	// err = godotenv.Load(envPath)
	// if err != nil {
	// 	log.Println("No .env file found, using environment variables")
	// }
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173" // Or your production frontend URL
	}

	// local testing
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")

	// Initialize database connection
	// database.InitDB(dbUser, dbPassword, dbName, dbHost, dbPort) // Pass dbPort here

	// Test on production state
	// Initialize database connection (using DATABASE_URL)
    database.InitDB() // No need to pass individual credentials anymore
	defer database.DB.Close()

	// Create a new router
	r := mux.NewRouter()

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

	r.HandleFunc("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserProfile(w, r, database.DB)
	}).Methods("GET")

	// Protected Routes (Login User Only)
	// User routes
	r.HandleFunc("/api/me", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCurrentUser(w, r, database.DB)
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

	r.HandleFunc("/api/upload-images", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UploadImages(w, r, database.DB)
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

	r.HandleFunc("/api/comments/{id}/upvote", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpvoteComment(w, r, database.DB)
	})).Methods("POST")

	r.HandleFunc("/api/comments/{id}/downvote", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.DownvoteComment(w, r, database.DB)
	})).Methods("POST")

	r.HandleFunc("/api/threads/{id}/like", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.LikeThread(w, r, database.DB)
	})).Methods("POST")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	// Create HTTP server
	srv := &http.Server{
		Handler:      handler,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
				log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	// Start server
	log.Printf("Server started on port: %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
	}

	log.Println("Server gracefully stopped")
}



// package main

// import (
// 	"log"
// 	"net/http"
// 	"travel-forum-backend/database"
// 	"travel-forum-backend/handlers"
// 	"travel-forum-backend/middleware"

// 	"github.com/gorilla/mux"
// 	"github.com/rs/cors"
// 	"github.com/joho/godotenv"
// )

// // func enableCORS(next http.Handler) http.Handler {
// // 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// // 		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174") // Replace with your React app's URL
// // 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// // 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// // 		if r.Method == "OPTIONS" {
// // 			return
// // 		}
// // 		next.ServeHTTP(w, r)
// // 	})
// // }

// func main() {
// 	// Initialize the database
// 	database.InitDB()
// 	defer database.DB.Close()

// 	// Create a new router
// 	r := mux.NewRouter()
// 	// r.Use(enableCORS)

	// // Public routes (Any user)
	// r.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.Register(w, r, database.DB)
	// }).Methods("POST")

	// r.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.Login(w, r, database.DB)
	// }).Methods("POST")

	// r.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.GetCategories(w, r, database.DB)
	// }).Methods("GET")

	// r.HandleFunc("/api/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.GetCategory(w, r, database.DB)
	// }).Methods("GET")

	// r.HandleFunc("/api/threads", func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.GetAllThreads(w, r, database.DB)
	// }).Methods("GET")

	// r.HandleFunc("/api/threads/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.GetThread(w, r, database.DB)
	// }).Methods("GET")

	// r.HandleFunc("/api/threads/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.GetCommentsByThread(w, r, database.DB)
	// }).Methods("GET")

	// r.HandleFunc("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.GetUserProfile(w, r, database.DB)
	// }).Methods("GET")

	// // Protected Routes (Login User Only)
	// // User routes
	// r.HandleFunc("/api/me", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.GetCurrentUser(w, r, database.DB)
	// })).Methods("GET")

	// r.HandleFunc("/api/users/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.UpdateUserProfile(w, r, database.DB)
	// })).Methods("PUT")

	// r.HandleFunc("/api/users/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.DeleteUser(w, r, database.DB)
	// })).Methods("DELETE")

	// // Categories routes
	// r.HandleFunc("/api/categories", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.CreateCategory(w, r, database.DB)
	// })).Methods("POST")

	// r.HandleFunc("/api/categories/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.UpdateCategory(w, r, database.DB)
	// })).Methods("PUT")

	// r.HandleFunc("/api/categories/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.DeleteCategory(w, r, database.DB)
	// })).Methods("DELETE")

	// // Thread routes
	// r.HandleFunc("/api/threads", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.CreateThread(w, r, database.DB)
	// })).Methods("POST")

	// r.HandleFunc("/api/upload-images", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.UploadImages(w, r, database.DB)
	// })).Methods("POST")

	// r.HandleFunc("/api/my-threads", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.GetMyThreads(w, r, database.DB)
	// }))

	// r.HandleFunc("/api/threads/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.UpdateThread(w, r, database.DB)
	// })).Methods("PUT")

	// r.HandleFunc("/api/threads/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.DeleteThread(w, r, database.DB)
	// })).Methods("DELETE")

	// // Comment routes
	// r.HandleFunc("/api/comments", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.CreateComment(w, r, database.DB)
	// })).Methods("POST")

	// r.HandleFunc("/api/comments/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.UpdateComment(w, r, database.DB)
	// })).Methods("PUT")

	// r.HandleFunc("/api/comments/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.DeleteComment(w, r, database.DB)
	// })).Methods("DELETE")

	// r.HandleFunc("/api/comments/{id}/upvote", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.UpvoteComment(w, r, database.DB)
	// })).Methods("POST")

	// r.HandleFunc("/api/comments/{id}/downvote", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.DownvoteComment(w, r, database.DB)
	// })).Methods("POST")

	// r.HandleFunc("/api/threads/{id}/like", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.LikeThread(w, r, database.DB)
	// })).Methods("POST")

// 	// Configure CORS
// 	c := cors.New(cors.Options{
// 		AllowedOrigins:   []string{"http://localhost:5173"}, // Allow your React frontend origin
// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowedHeaders:   []string{"Content-Type", "Authorization"},
// 		AllowCredentials: true,
// 	})

// 	// Wrap your router with the CORS middleware
// 	handler := c.Handler(r)

// 	// Start the server
// 	log.Println("Server started on :8080")
// 	log.Fatal(http.ListenAndServe(":8080", handler))
// }

