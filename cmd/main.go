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
		// Algorithm information endpoint
		api.GET("/algorithms", optimizationHandler.GetSupportedAlgorithms)

		// Money change algorithm (used for cash payments)
		api.POST("/change", optimizationHandler.CalculateChange)

		// Inventory optimization (Knapsack algorithm)
		api.POST("/inventory/optimize", optimizationHandler.OptimizeInventory)

		// Table assignment optimization (Greedy algorithm)
		api.POST("/tables/assign", optimizationHandler.AssignTables)
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
