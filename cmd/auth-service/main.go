package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sashakosti/auth-service/internal/token"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Config holds application configuration
type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

// App represents the application
type App struct {
	DB           *sql.DB
	TokenService token.TokenManager
	Config       *Config
}

// AuthRequest represents login/refresh request body
type AuthRequest struct {
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using environment variables")
	}
}

func main() {
	// Load configuration
	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"),
		JWTSecret:       getEnv("JWT_SECRET", "default-secret-key-must-be-at-least-32-characters-long"),
		AccessTokenExp:  15 * time.Minute,
		RefreshTokenExp: 7 * 24 * time.Hour,
	}

	// Initialize database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify database connection
	if err = db.Ping(); err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}

	// Initialize token service
	tokenService := token.NewTokenManager(
		cfg.JWTSecret,
		cfg.AccessTokenExp,
		cfg.RefreshTokenExp,
		nil, // TODO You'll need to provide a proper storage implementation here
	)

	// Create app
	app := &App{
		DB:           db,
		TokenService: tokenService,
		Config:       cfg,
	}

	// Initialize router
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/api/health", app.healthCheck).Methods("GET")
	r.HandleFunc("/api/login", app.login).Methods("POST")
	r.HandleFunc("/api/refresh", app.refreshToken).Methods("POST")

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(app.authMiddleware)
	api.HandleFunc("/protected", app.protectedHandler).Methods("GET")

	// Start server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("🚀 Server started on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Health check endpoint
func (app *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Login handler
func (app *App) login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// In a real app, validate user credentials here
	if req.UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Generate tokens
	tokenPair, err := app.TokenService.GenerateTokens(req.UserID)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// In a real app, you might want to store the refresh token in the database here
	// For example: saveRefreshTokenToDB(req.UserID, tokenPair.RefreshToken)

	// Respond with tokens
	json.NewEncoder(w).Encode(AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		RefreshToken: tokenPair.RefreshToken,
	})
}

// Refresh token handler
func (app *App) refreshToken(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.RefreshToken == "" {
		http.Error(w, "User ID and refresh token are required", http.StatusBadRequest)
		return
	}

	// Use the TokenManager's RefreshTokens method which handles both validation and new token generation
	tokenPair, err := app.TokenService.RefreshTokens(req.UserID, req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// In a real app, you might want to update the stored refresh token here
	// For example: updateRefreshTokenInDB(req.UserID, tokenPair.RefreshToken)

	// Respond with new tokens
	json.NewEncoder(w).Encode(AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		RefreshToken: tokenPair.RefreshToken,
	})
}

// Auth middleware
func (app *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Extract the token from the Authorization header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]
		userID, err := app.TokenService.ValidateAccessToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user ID to request context
		r = r.WithContext(context.WithValue(r.Context(), "userID", userID))

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Protected endpoint example
func (app *App) protectedHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Hello, user %s! This is a protected endpoint.", userID),
	})
}
