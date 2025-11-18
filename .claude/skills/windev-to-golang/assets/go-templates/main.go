package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
	"go.uber.org/zap"
)

var (
	db     *sql.DB
	logger *zap.Logger
	config Config
)

// Config holds application configuration
type Config struct {
	DatabaseURL string
	Port        string
	Environment string
}

func main() {
	// Initialize logger
	var err error
	if os.Getenv("ENVIRONMENT") == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("starting application")

	// Load configuration
	config = LoadConfig()

	// Initialize database
	if err := InitDatabase(); err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	// Setup router
	router := SetupRouter()

	// Start server
	port := config.Port
	if port == "" {
		port = "8080"
	}

	logger.Info("server starting", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		Environment: os.Getenv("ENVIRONMENT"),
	}
}

// InitDatabase initializes database connection
func InitDatabase() error {
	var err error
	db, err = sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	logger.Info("database connection established")
	return nil
}

// SetupRouter configures all routes and middleware
func SetupRouter() *gin.Engine {
	// Set Gin mode
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(LoggingMiddleware())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// Health check endpoint
	router.GET("/health", HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// Add your routes here
		// Example:
		// api.GET("/users/:id", GetUser)
		// api.POST("/users", CreateUser)
	}

	return router
}

// LoggingMiddleware logs all HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		// start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Log after request
		statusCode := c.Writer.Status()
		// duration := time.Since(start)

		logger.Info("request completed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			// zap.Duration("duration", duration),
		)
	}
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// HealthCheck returns the health status of the service
func HealthCheck(c *gin.Context) {
	// Check database connection
	if err := db.Ping(); err != nil {
		logger.Error("health check failed", zap.Error(err))
		c.JSON(503, gin.H{
			"status":   "unhealthy",
			"database": "disconnected",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":   "healthy",
		"database": "connected",
	})
}
