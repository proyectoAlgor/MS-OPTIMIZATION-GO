package handlers

import (
	"ms-optimization-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OptimizationHandler handles HTTP requests for optimization algorithms
type OptimizationHandler struct {
	optimizationService *service.OptimizationService
}

// NewOptimizationHandler creates a new optimization handler
func NewOptimizationHandler() *OptimizationHandler {
	return &OptimizationHandler{
		optimizationService: service.NewOptimizationService(),
	}
}

// HealthCheck returns the health status of the service
func (h *OptimizationHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "ms-optimization-go",
		"status":  "healthy",
		"algorithms": []string{
			"money_change",
			"knapsack",
			"table_assignment",
		},
	})
}

// CalculateChange handles change calculation requests
func (h *OptimizationHandler) CalculateChange(c *gin.Context) {
	var req service.CalculateChangeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if req.AmountPaid < 0 || req.TotalCost < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Amount paid and total cost must be non-negative",
		})
		return
	}

	result := h.optimizationService.CalculateOptimalChange(req)

	status := http.StatusOK
	if !result.Success {
		status = http.StatusBadRequest
	}

	c.JSON(status, result)
}

// GetSupportedAlgorithms returns information about supported algorithms
func (h *OptimizationHandler) GetSupportedAlgorithms(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"algorithms": gin.H{
			"money_change": gin.H{
				"description": "Greedy algorithm for optimal coin change",
				"complexity":  "O(n log n) for sorting + O(n) for processing",
				"use_case":    "Calculate optimal change when customer pays in cash",
			},
			"knapsack": gin.H{
				"description": "0/1 Knapsack and Fractional Knapsack algorithms for inventory optimization",
				"algorithms":  []string{"dynamic_programming", "greedy"},
				"complexity": gin.H{
					"dynamic_programming": "O(n * capacity) - exact solution",
					"greedy":              "O(n log n) - approximate solution, faster",
				},
				"use_case": "Optimize which products to keep in stock considering space, demand, and cost",
			},
			"table_assignment": gin.H{
				"description": "Greedy algorithm for optimal table assignment to customer groups",
				"algorithms":  []string{"greedy", "optimal"},
				"complexity": gin.H{
					"greedy":  "O(n * m) where n = groups, m = tables - fast and efficient",
					"optimal": "O(n³) - optimal solution but slower",
				},
				"use_case": "Assign tables to customers considering capacity, proximity, and priority",
			},
		},
		"message": "Supported optimization algorithms for bar management",
	})
}

// OptimizeInventory handles inventory optimization requests using Knapsack algorithm
func (h *OptimizationHandler) OptimizeInventory(c *gin.Context) {
	var req service.OptimizeInventoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if req.MaxCapacity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Max capacity must be greater than 0",
		})
		return
	}

	// Validate algorithm
	if req.Algorithm != "" && req.Algorithm != "dp" && req.Algorithm != "greedy" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":       false,
			"error":         "Invalid algorithm",
			"valid_options": []string{"dp", "greedy"},
		})
		return
	}

	result := h.optimizationService.OptimizeInventory(req)

	status := http.StatusOK
	if !result.Success {
		status = http.StatusBadRequest
	}

	c.JSON(status, result)
}

// AssignTables handles table assignment optimization requests
func (h *OptimizationHandler) AssignTables(c *gin.Context) {
	var req service.AssignTablesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if len(req.Tables) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No tables provided",
		})
		return
	}

	if len(req.Groups) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No customer groups provided",
		})
		return
	}

	// Validate method
	if req.Method != "" && req.Method != "greedy" && req.Method != "optimal" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":       false,
			"error":         "Invalid method",
			"valid_options": []string{"greedy", "optimal"},
		})
		return
	}

	result := h.optimizationService.AssignTables(req)

	status := http.StatusOK
	if !result.Success {
		status = http.StatusBadRequest
	}

	c.JSON(status, result)
}
