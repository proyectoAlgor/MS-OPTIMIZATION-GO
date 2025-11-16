package main

import (
	"ms-optimization-go/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get environment variables
	port := getEnv("PORT", "8080")

	// Initialize handler
	optimizationHandler := handlers.NewOptimizationHandler()

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", optimizationHandler.HealthCheck)

	// API routes
	api := r.Group("/api/optimization")
	{
		// Algorithm information endpoints
		api.GET("/coins", optimizationHandler.GetAvailableCoins)
		api.GET("/algorithms", optimizationHandler.GetSupportedAlgorithms)

		// Money change algorithm
		api.POST("/change", optimizationHandler.CalculateChange)

		// Sorting algorithms
		api.POST("/sort/products", optimizationHandler.SortProducts)

		// Search algorithms
		api.POST("/search/products", optimizationHandler.SearchProducts)

		// Order analysis
		api.POST("/analyze/order", optimizationHandler.AnalyzeOrder)
	}

	// Start server
	gin.SetMode(gin.ReleaseMode)
	r.Run(":" + port)
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
