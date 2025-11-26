package service

import (
	"fmt"
	"ms-optimization-go/internal/algorithms"
)

// OptimizationService provides business logic for optimization algorithms
type OptimizationService struct {
	moneyAlgo       *algorithms.MoneyChangeAlgorithm
	knapsackAlgo    *algorithms.KnapsackAlgorithm
	tableAssignAlgo *algorithms.TableAssignmentAlgorithm
}

// NewOptimizationService creates a new optimization service
func NewOptimizationService() *OptimizationService {
	// Initialize with common coin denominations (in cents)
	coins := []int{5000, 2000, 1000, 500, 200, 100, 50, 25, 10, 5, 1} // $50, $20, $10, $5, $2, $1, $0.50, $0.25, $0.10, $0.05, $0.01

	return &OptimizationService{
		moneyAlgo:       algorithms.NewMoneyChangeAlgorithm(coins),
		knapsackAlgo:    algorithms.NewKnapsackAlgorithm(),
		tableAssignAlgo: algorithms.NewTableAssignmentAlgorithm(),
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

// OptimizeInventoryRequest represents a request to optimize inventory
type OptimizeInventoryRequest struct {
	Items          []algorithms.InventoryItem `json:"items"`
	MaxCapacity    float64                    `json:"max_capacity"`
	MinDemandScore float64                    `json:"min_demand_score,omitempty"`
	Algorithm      string                     `json:"algorithm"` // "dp" (dynamic programming) or "greedy"
}

// OptimizeInventoryResponse represents the response for inventory optimization
type OptimizeInventoryResponse struct {
	Success bool                      `json:"success"`
	Result  algorithms.KnapsackResult `json:"result"`
	Message string                    `json:"message"`
}

// OptimizeInventory optimizes inventory selection using Knapsack algorithm
func (os *OptimizationService) OptimizeInventory(req OptimizeInventoryRequest) OptimizeInventoryResponse {
	if len(req.Items) == 0 {
		return OptimizeInventoryResponse{
			Success: false,
			Message: "No items provided",
		}
	}

	if req.MaxCapacity <= 0 {
		return OptimizeInventoryResponse{
			Success: false,
			Message: "Invalid capacity (must be > 0)",
		}
	}

	minDemandScore := req.MinDemandScore
	if minDemandScore < 0 {
		minDemandScore = 0
	}

	var result algorithms.KnapsackResult
	if req.Algorithm == "dp" {
		// Use dynamic programming for exact solution
		result = os.knapsackAlgo.SolveKnapsack01(req.Items, req.MaxCapacity)
	} else if req.Algorithm == "greedy" || req.Algorithm == "" {
		// Use greedy approach for faster results
		if minDemandScore > 0 {
			result = os.knapsackAlgo.OptimizeInventory(req.Items, req.MaxCapacity, minDemandScore)
		} else {
			result = os.knapsackAlgo.SolveKnapsackGreedy(req.Items, req.MaxCapacity)
		}
	} else {
		return OptimizeInventoryResponse{
			Success: false,
			Message: "Invalid algorithm. Supported: 'dp' (dynamic programming) or 'greedy'",
		}
	}

	return OptimizeInventoryResponse{
		Success: true,
		Result:  result,
		Message: fmt.Sprintf("Inventory optimized: %d items selected", len(result.SelectedItems)),
	}
}

// AssignTablesRequest represents a request to assign tables to customer groups
type AssignTablesRequest struct {
	Tables []algorithms.TableAssignment `json:"tables"`
	Groups []algorithms.CustomerGroup   `json:"groups"`
	Method string                       `json:"method"` // "greedy" or "optimal"
}

// AssignTablesResponse represents the response for table assignment
type AssignTablesResponse struct {
	Success bool                        `json:"success"`
	Result  algorithms.AssignmentResult `json:"result"`
	Message string                      `json:"message"`
}

// AssignTables assigns tables to customer groups using optimization algorithms
func (os *OptimizationService) AssignTables(req AssignTablesRequest) AssignTablesResponse {
	if len(req.Tables) == 0 {
		return AssignTablesResponse{
			Success: false,
			Message: "No tables provided",
		}
	}

	if len(req.Groups) == 0 {
		return AssignTablesResponse{
			Success: false,
			Message: "No customer groups provided",
		}
	}

	var result algorithms.AssignmentResult
	if req.Method == "optimal" {
		result = os.tableAssignAlgo.AssignTablesOptimal(req.Tables, req.Groups)
	} else {
		// Default to greedy (faster and usually good enough)
		result = os.tableAssignAlgo.AssignTablesGreedy(req.Tables, req.Groups)
	}

	return AssignTablesResponse{
		Success: true,
		Result:  result,
		Message: result.Message,
	}
}
