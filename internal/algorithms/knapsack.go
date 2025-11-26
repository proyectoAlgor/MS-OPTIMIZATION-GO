package algorithms

import (
	"fmt"
	"sort"
)

// KnapsackAlgorithm implements the 0/1 Knapsack problem for inventory optimization
// This is useful for deciding which products to keep in stock considering space, demand, and cost
type KnapsackAlgorithm struct{}

// NewKnapsackAlgorithm creates a new instance
func NewKnapsackAlgorithm() *KnapsackAlgorithm {
	return &KnapsackAlgorithm{}
}

// InventoryItem represents an item that can be stored in inventory
type InventoryItem struct {
	ID          string  // Product ID
	Name        string  // Product name
	Weight      float64 // Space/weight occupied (e.g., cubic meters, kg)
	Value       float64 // Value/priority (e.g., demand score, profit margin)
	Cost        float64 // Cost of the item
	DemandScore float64 // Demand score (0-1, where 1 is highest demand)
}

// KnapsackResult represents the result of the knapsack optimization
type KnapsackResult struct {
	SelectedItems     []InventoryItem `json:"selected_items"`
	TotalValue        float64         `json:"total_value"`
	TotalWeight       float64         `json:"total_weight"`
	TotalCost         float64         `json:"total_cost"`
	CapacityUsed      float64         `json:"capacity_used"`
	CapacityAvailable float64         `json:"capacity_available"`
	Efficiency        float64         `json:"efficiency"` // value/weight ratio
	Message           string          `json:"message"`
}

// SolveKnapsack01 solves the 0/1 Knapsack problem using dynamic programming
// Time Complexity: O(n * capacity)
// Space Complexity: O(n * capacity)
// This is optimal for exact solutions but can be slow for large inputs
func (ka *KnapsackAlgorithm) SolveKnapsack01(items []InventoryItem, capacity float64) KnapsackResult {
	if len(items) == 0 {
		return KnapsackResult{
			SelectedItems:     []InventoryItem{},
			CapacityAvailable: capacity,
			Message:           "No items provided",
		}
	}

	if capacity <= 0 {
		return KnapsackResult{
			SelectedItems:     []InventoryItem{},
			CapacityAvailable: capacity,
			Message:           "Invalid capacity (must be > 0)",
		}
	}

	// Convert capacity to integer for DP (multiply by 100 to preserve 2 decimal places)
	capacityInt := int(capacity * 100)
	n := len(items)

	// DP table: dp[i][w] = maximum value with first i items and weight w
	// We use 1D array to optimize space: dp[w] = maximum value with weight w
	dp := make([]float64, capacityInt+1)
	selected := make([][]bool, n)
	for i := range selected {
		selected[i] = make([]bool, capacityInt+1)
	}

	// Build DP table
	for i := 0; i < n; i++ {
		weightInt := int(items[i].Weight * 100)
		value := items[i].Value

		// Iterate backwards to avoid using the same item twice
		for w := capacityInt; w >= weightInt; w-- {
			if dp[w-weightInt]+value > dp[w] {
				dp[w] = dp[w-weightInt] + value
				selected[i][w] = true
			}
		}
	}

	// Reconstruct the solution
	selectedItems := []InventoryItem{}
	totalWeight := 0.0
	totalValue := 0.0
	totalCost := 0.0

	w := capacityInt
	for i := n - 1; i >= 0; i-- {
		if w >= int(items[i].Weight*100) && selected[i][w] {
			selectedItems = append(selectedItems, items[i])
			totalWeight += items[i].Weight
			totalValue += items[i].Value
			totalCost += items[i].Cost
			w -= int(items[i].Weight * 100)
		}
	}

	// Reverse to maintain original order
	for i, j := 0, len(selectedItems)-1; i < j; i, j = i+1, j-1 {
		selectedItems[i], selectedItems[j] = selectedItems[j], selectedItems[i]
	}

	efficiency := 0.0
	if totalWeight > 0 {
		efficiency = totalValue / totalWeight
	}

	return KnapsackResult{
		SelectedItems:     selectedItems,
		TotalValue:        totalValue,
		TotalWeight:       totalWeight,
		TotalCost:         totalCost,
		CapacityUsed:      totalWeight,
		CapacityAvailable: capacity - totalWeight,
		Efficiency:        efficiency,
		Message:           fmt.Sprintf("Selected %d items with total value %.2f and weight %.2f/%.2f", len(selectedItems), totalValue, totalWeight, capacity),
	}
}

// SolveKnapsackGreedy solves the fractional knapsack problem using greedy approach
// Time Complexity: O(n log n) for sorting + O(n) for processing
// Space Complexity: O(n)
// This is faster but may not be optimal for 0/1 knapsack (though often very close)
// Best used when items can be partially included (fractional knapsack)
func (ka *KnapsackAlgorithm) SolveKnapsackGreedy(items []InventoryItem, capacity float64) KnapsackResult {
	if len(items) == 0 {
		return KnapsackResult{
			SelectedItems:     []InventoryItem{},
			CapacityAvailable: capacity,
			Message:           "No items provided",
		}
	}

	if capacity <= 0 {
		return KnapsackResult{
			SelectedItems:     []InventoryItem{},
			CapacityAvailable: capacity,
			Message:           "Invalid capacity (must be > 0)",
		}
	}

	// Create a copy of items with their value/weight ratio
	type ItemWithRatio struct {
		Item          InventoryItem
		Ratio         float64
		OriginalIndex int
	}

	itemsWithRatio := make([]ItemWithRatio, len(items))
	for i, item := range items {
		ratio := 0.0
		if item.Weight > 0 {
			ratio = item.Value / item.Weight
		}
		itemsWithRatio[i] = ItemWithRatio{
			Item:          item,
			Ratio:         ratio,
			OriginalIndex: i,
		}
	}

	// Sort by value/weight ratio (descending) - greedy approach
	sort.Slice(itemsWithRatio, func(i, j int) bool {
		return itemsWithRatio[i].Ratio > itemsWithRatio[j].Ratio
	})

	selectedItems := []InventoryItem{}
	remainingCapacity := capacity
	totalValue := 0.0
	totalWeight := 0.0
	totalCost := 0.0

	// Greedy selection: take items with highest value/weight ratio first
	for _, itemWithRatio := range itemsWithRatio {
		if remainingCapacity <= 0 {
			break
		}

		item := itemWithRatio.Item
		if item.Weight <= remainingCapacity {
			// Take the whole item
			selectedItems = append(selectedItems, item)
			totalValue += item.Value
			totalWeight += item.Weight
			totalCost += item.Cost
			remainingCapacity -= item.Weight
		} else if remainingCapacity > 0 {
			// Take a fraction of the item (for fractional knapsack)
			// For 0/1 knapsack, we skip it, but for demonstration we include it partially
			fraction := remainingCapacity / item.Weight
			partialItem := item
			partialItem.Weight = remainingCapacity
			partialItem.Value = item.Value * fraction
			partialItem.Cost = item.Cost * fraction

			selectedItems = append(selectedItems, partialItem)
			totalValue += partialItem.Value
			totalWeight += remainingCapacity
			totalCost += partialItem.Cost
			remainingCapacity = 0
			break
		}
	}

	efficiency := 0.0
	if totalWeight > 0 {
		efficiency = totalValue / totalWeight
	}

	return KnapsackResult{
		SelectedItems:     selectedItems,
		TotalValue:        totalValue,
		TotalWeight:       totalWeight,
		TotalCost:         totalCost,
		CapacityUsed:      totalWeight,
		CapacityAvailable: remainingCapacity,
		Efficiency:        efficiency,
		Message:           fmt.Sprintf("Selected %d items (greedy) with total value %.2f and weight %.2f/%.2f", len(selectedItems), totalValue, totalWeight, capacity),
	}
}

// OptimizeInventory optimizes inventory selection based on demand, space, and cost
// This is a specialized version for bar inventory management
func (ka *KnapsackAlgorithm) OptimizeInventory(items []InventoryItem, maxCapacity float64, minDemandScore float64) KnapsackResult {
	// Filter items by minimum demand score
	filteredItems := []InventoryItem{}
	for _, item := range items {
		if item.DemandScore >= minDemandScore {
			filteredItems = append(filteredItems, item)
		}
	}

	if len(filteredItems) == 0 {
		return KnapsackResult{
			SelectedItems:     []InventoryItem{},
			CapacityAvailable: maxCapacity,
			Message:           fmt.Sprintf("No items meet the minimum demand score of %.2f", minDemandScore),
		}
	}

	// Use greedy approach for faster results (good enough for inventory management)
	return ka.SolveKnapsackGreedy(filteredItems, maxCapacity)
}
