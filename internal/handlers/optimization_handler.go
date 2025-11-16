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
			"sorting",
			"search",
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

// SortProducts handles product sorting requests
func (h *OptimizationHandler) SortProducts(c *gin.Context) {
	var req service.SortProductsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate sort criteria
	validSortBy := map[string]bool{
		"price_asc":    true,
		"price_desc":   true,
		"name_asc":     true,
		"name_desc":    true,
		"code_asc":     true,
		"category_asc": true,
	}

	if !validSortBy[req.SortBy] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":       false,
			"error":         "Invalid sort criteria",
			"valid_options": []string{"price_asc", "price_desc", "name_asc", "name_desc", "code_asc", "category_asc"},
		})
		return
	}

	// Validate algorithm
	validAlgorithms := map[string]bool{
		"quick":     true,
		"insertion": true,
		"selection": true,
	}

	if !validAlgorithms[req.Algorithm] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":       false,
			"error":         "Invalid algorithm",
			"valid_options": []string{"quick", "insertion", "selection"},
		})
		return
	}

	result := h.optimizationService.SortProducts(req)

	status := http.StatusOK
	if !result.Success {
		status = http.StatusBadRequest
	}

	c.JSON(status, result)
}

// SearchProducts handles product search requests
func (h *OptimizationHandler) SearchProducts(c *gin.Context) {
	var req service.SearchProductsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate search type
	validSearchTypes := map[string]bool{
		"name":        true,
		"code":        true,
		"price_range": true,
		"price_exact": true,
	}

	if !validSearchTypes[req.SearchType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":       false,
			"error":         "Invalid search type",
			"valid_options": []string{"name", "code", "price_range", "price_exact"},
		})
		return
	}

	result := h.optimizationService.SearchProducts(req)

	status := http.StatusOK
	if !result.Success {
		status = http.StatusBadRequest
	}

	c.JSON(status, result)
}

// AnalyzeOrder handles order analysis requests
func (h *OptimizationHandler) AnalyzeOrder(c *gin.Context) {
	var req service.AnalyzeOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	result := h.optimizationService.AnalyzeOrder(req)

	status := http.StatusOK
	if !result.Success {
		status = http.StatusBadRequest
	}

	c.JSON(status, result)
}

// GetAvailableCoins returns the available coin denominations
func (h *OptimizationHandler) GetAvailableCoins(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"coins":   []string{"$50.00", "$20.00", "$10.00", "$5.00", "$2.00", "$1.00", "$0.50", "$0.25", "$0.10", "$0.05", "$0.01"},
		"message": "Available coin denominations for change calculation",
	})
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
			"sorting": gin.H{
				"description": "Various sorting algorithms for products and data",
				"algorithms":  []string{"quick_sort", "insertion_sort", "selection_sort"},
				"complexity": gin.H{
					"quick_sort":     "O(n log n) average, O(n²) worst case",
					"insertion_sort": "O(n²) average, O(n) best case",
					"selection_sort": "O(n²) in all cases",
				},
				"use_case": "Sort products by price, name, category, etc.",
			},
			"search": gin.H{
				"description": "Search algorithms for finding products and data",
				"algorithms":  []string{"binary_search", "linear_search", "string_reversal"},
				"complexity": gin.H{
					"binary_search":   "O(log n)",
					"linear_search":   "O(n)",
					"string_reversal": "O(n)",
				},
				"use_case": "Find products by name, code, price range",
			},
		},
		"message": "Supported optimization algorithms for bar management",
	})
}
