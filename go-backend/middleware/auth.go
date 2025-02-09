package middleware

import (
	"context"
	"log"
	"net/http"
	// "os"
	// "path/filepath"

	// "os"
	"strings"

	"travel-forum-backend/models"

	"github.com/golang-jwt/jwt"
	// "github.com/joho/godotenv"
)

// func init() {
// 	// Construct path to .env file in the root directory
// 	rootDir, err := os.Getwd()
// 	if err != nil {
// 		log.Fatal("Error getting working directory:", err)
// 	}

// 	envPath := filepath.Join(rootDir, "..", ".env") // Go up one level for .env
// 	err = godotenv.Load(envPath)
// 	if err != nil {
// 		log.Println("No .env file found, using environment variables")
// 	}
// }

// var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var jwtKey = []byte("5lfX8Bl4C1mZZ/ljU+BrWFoxTcxQqacwPVfloDs+5No=")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic in AuthMiddleware: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		if len(jwtKey) == 0 {
			log.Println("JWT secret key not set")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Get token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Authorization header missing")
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Extract token from header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			log.Println("Invalid token format")
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Parse the token
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.Println("Invalid token signature")
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					log.Println("Token expired")
					http.Error(w, "Token expired", http.StatusUnauthorized)
					return
				}
			}
			log.Printf("Invalid token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			log.Println("Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Validate token claims
		if claims.UserID == 0 {
			log.Println("Invalid token claims: UserID missing")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add the claims to the request context
		ctx := context.WithValue(r.Context(), "claims", claims)
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		r = r.WithContext(ctx)

		log.Printf("Token validated successfully. UserID: %d", claims.UserID)

		// Call the next middleware or handler
		next(w, r)
	}
}
