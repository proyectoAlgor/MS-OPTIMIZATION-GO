package service

import (
	"fmt"
	"ms-optimization-go/internal/algorithms"
)

// OptimizationService provides business logic for optimization algorithms
type OptimizationService struct {
	moneyAlgo   *algorithms.MoneyChangeAlgorithm
	sortingAlgo *algorithms.SortingAlgorithm
	searchAlgo  *algorithms.SearchAlgorithm
}

// NewOptimizationService creates a new optimization service
func NewOptimizationService() *OptimizationService {
	// Initialize with common coin denominations (in cents)
	coins := []int{5000, 2000, 1000, 500, 200, 100, 50, 25, 10, 5, 1} // $50, $20, $10, $5, $2, $1, $0.50, $0.25, $0.10, $0.05, $0.01

	return &OptimizationService{
		moneyAlgo:   algorithms.NewMoneyChangeAlgorithm(coins),
		sortingAlgo: algorithms.NewSortingAlgorithm(),
		searchAlgo:  algorithms.NewSearchAlgorithm(),
	}
}

// CalculateChangeRequest represents a request to calculate change
type CalculateChangeRequest struct {
	AmountPaid float64 `json:"amount_paid"`
	TotalCost  float64 `json:"total_cost"`
}

// CalculateChangeResponse represents the response for change calculation
type CalculateChangeResponse struct {
	Success        bool           `json:"success"`
	ChangeAmount   float64        `json:"change_amount"`
	TotalCoins     int            `json:"total_coins"`
	Breakdown      map[string]int `json:"breakdown"`
	Message        string         `json:"message"`
	AvailableCoins []string       `json:"available_coins"`
}

// CalculateOptimalChange calculates the optimal change for a payment
func (os *OptimizationService) CalculateOptimalChange(req CalculateChangeRequest) CalculateChangeResponse {
	// Convert to cents to avoid floating point precision issues
	amountPaidCents := int(req.AmountPaid * 100)
	totalCostCents := int(req.TotalCost * 100)

	changeAmountCents := amountPaidCents - totalCostCents

	if changeAmountCents < 0 {
		return CalculateChangeResponse{
			Success:      false,
			ChangeAmount: 0,
			Message:      "Insufficient payment amount",
		}
	}

	if changeAmountCents == 0 {
		return CalculateChangeResponse{
			Success:        true,
			ChangeAmount:   0,
			TotalCoins:     0,
			Breakdown:      make(map[string]int),
			Message:        "Exact payment, no change needed",
			AvailableCoins: os.formatCoins(os.moneyAlgo.GetAvailableCoins()),
		}
	}

	result := os.moneyAlgo.CalculateChange(changeAmountCents)

	// Convert breakdown from cents to dollar format
	breakdown := make(map[string]int)
	for coinValue, quantity := range result.Breakdown {
		dollarValue := float64(coinValue) / 100
		breakdown[fmt.Sprintf("$%.2f", dollarValue)] = quantity
	}

	return CalculateChangeResponse{
		Success:        result.Success,
		ChangeAmount:   float64(changeAmountCents) / 100,
		TotalCoins:     result.TotalCoins,
		Breakdown:      breakdown,
		Message:        result.Message,
		AvailableCoins: os.formatCoins(os.moneyAlgo.GetAvailableCoins()),
	}
}

// formatCoins formats coin values from cents to dollar format
func (os *OptimizationService) formatCoins(coins []int) []string {
	formatted := make([]string, len(coins))
	for i, coin := range coins {
		dollarValue := float64(coin) / 100
		formatted[i] = fmt.Sprintf("$%.2f", dollarValue)
	}
	return formatted
}

// SortProductsRequest represents a request to sort products
type SortProductsRequest struct {
	Products  []algorithms.Product `json:"products"`
	SortBy    string               `json:"sort_by"`   // price_asc, price_desc, name_asc, name_desc, code_asc, category_asc
	Algorithm string               `json:"algorithm"` // quick, insertion, selection
}

// SortProductsResponse represents the response for sorting products
type SortProductsResponse struct {
	Success   bool                 `json:"success"`
	Products  []algorithms.Product `json:"products"`
	Message   string               `json:"message"`
	Algorithm string               `json:"algorithm_used"`
}

// SortProducts sorts products using the specified algorithm
func (os *OptimizationService) SortProducts(req SortProductsRequest) SortProductsResponse {
	if len(req.Products) == 0 {
		return SortProductsResponse{
			Success:   false,
			Message:   "No products provided",
			Algorithm: req.Algorithm,
		}
	}

	var sortedProducts []algorithms.Product
	var message string

	switch req.Algorithm {
	case "quick":
		sortedProducts = os.sortingAlgo.QuickSortProducts(req.Products, req.SortBy)
		message = "Products sorted using Quick Sort algorithm"
	case "insertion":
		sortedProducts = os.sortingAlgo.InsertionSortProducts(req.Products, req.SortBy)
		message = "Products sorted using Insertion Sort algorithm (optimal for small lists)"
	case "selection":
		sortedProducts = os.sortingAlgo.SelectionSortProducts(req.Products, req.SortBy)
		message = "Products sorted using Selection Sort algorithm"
	default:
		// Default to quick sort
		sortedProducts = os.sortingAlgo.QuickSortProducts(req.Products, req.SortBy)
		message = "Products sorted using Quick Sort algorithm (default)"
	}

	return SortProductsResponse{
		Success:   true,
		Products:  sortedProducts,
		Message:   message,
		Algorithm: req.Algorithm,
	}
}

// SearchProductsRequest represents a request to search products
type SearchProductsRequest struct {
	Products   []algorithms.Product `json:"products"`
	SearchType string               `json:"search_type"` // name, code, price_range, price_exact
	SearchTerm string               `json:"search_term"`
	MinPrice   *float64             `json:"min_price,omitempty"`
	MaxPrice   *float64             `json:"max_price,omitempty"`
	ExactPrice *float64             `json:"exact_price,omitempty"`
}

// SearchProductsResponse represents the response for searching products
type SearchProductsResponse struct {
	Success  bool                 `json:"success"`
	Products []algorithms.Product `json:"products"`
	Message  string               `json:"message"`
	Total    float64              `json:"total_value,omitempty"`
}

// SearchProducts searches for products using various algorithms
func (os *OptimizationService) SearchProducts(req SearchProductsRequest) SearchProductsResponse {
	if len(req.Products) == 0 {
		return SearchProductsResponse{
			Success: false,
			Message: "No products provided",
		}
	}

	var result []algorithms.Product
	var message string

	switch req.SearchType {
	case "name":
		result = os.searchAlgo.SearchProductsByName(req.Products, req.SearchTerm)
		message = fmt.Sprintf("Found %d products matching name '%s'", len(result), req.SearchTerm)
	case "code":
		product := os.searchAlgo.SearchProductsByCode(req.Products, req.SearchTerm)
		if product != nil {
			result = []algorithms.Product{*product}
			message = fmt.Sprintf("Found product with code '%s'", req.SearchTerm)
		} else {
			result = []algorithms.Product{}
			message = fmt.Sprintf("No product found with code '%s'", req.SearchTerm)
		}
	case "price_range":
		if req.MinPrice != nil && req.MaxPrice != nil {
			result = os.searchAlgo.BinarySearchProductsByPriceRange(req.Products, *req.MinPrice, *req.MaxPrice)
			message = fmt.Sprintf("Found %d products in price range $%.2f - $%.2f", len(result), *req.MinPrice, *req.MaxPrice)
		} else {
			return SearchProductsResponse{
				Success: false,
				Message: "MinPrice and MaxPrice are required for price range search",
			}
		}
	case "price_exact":
		if req.ExactPrice != nil {
			searchResult := os.searchAlgo.BinarySearchProducts(req.Products, *req.ExactPrice)
			if searchResult.Found {
				// We need to get the actual product from the original list
				// Since binary search returns the index from sorted list
				message = searchResult.Message
			} else {
				result = []algorithms.Product{}
				message = searchResult.Message
			}
		} else {
			return SearchProductsResponse{
				Success: false,
				Message: "ExactPrice is required for exact price search",
			}
		}
	default:
		return SearchProductsResponse{
			Success: false,
			Message: "Invalid search type. Supported types: name, code, price_range, price_exact",
		}
	}

	total := os.searchAlgo.SumProductPrices(result)

	return SearchProductsResponse{
		Success:  true,
		Products: result,
		Message:  message,
		Total:    total,
	}
}

// AnalyzeOrderRequest represents a request to analyze an order
type AnalyzeOrderRequest struct {
	Products []algorithms.Product `json:"products"`
}

// AnalyzeOrderResponse represents the response for order analysis
type AnalyzeOrderResponse struct {
	Success        bool                `json:"success"`
	Total          float64             `json:"total"`
	TotalRecursive float64             `json:"total_recursive"`
	ProductCount   int                 `json:"product_count"`
	MostExpensive  *algorithms.Product `json:"most_expensive_product,omitempty"`
	Cheapest       *algorithms.Product `json:"cheapest_product,omitempty"`
	Message        string              `json:"message"`
}

// AnalyzeOrder analyzes an order using various algorithms
func (os *OptimizationService) AnalyzeOrder(req AnalyzeOrderRequest) AnalyzeOrderResponse {
	if len(req.Products) == 0 {
		return AnalyzeOrderResponse{
			Success: false,
			Message: "No products in order",
		}
	}

	total := os.searchAlgo.SumProductPrices(req.Products)
	totalRecursive := os.searchAlgo.SumProductPricesRecursive(req.Products)
	mostExpensive := os.searchAlgo.FindMostExpensiveProduct(req.Products)
	cheapest := os.searchAlgo.FindCheapestProduct(req.Products)

	return AnalyzeOrderResponse{
		Success:        true,
		Total:          total,
		TotalRecursive: totalRecursive,
		ProductCount:   len(req.Products),
		MostExpensive:  mostExpensive,
		Cheapest:       cheapest,
		Message:        fmt.Sprintf("Order analyzed: %d products, total $%.2f", len(req.Products), total),
	}
}
