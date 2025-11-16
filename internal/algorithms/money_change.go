package algorithms

import (
	"fmt"
	"sort"
)

// MoneyChangeAlgorithm implements a greedy algorithm for optimal coin change
type MoneyChangeAlgorithm struct {
	coins []int
}

// NewMoneyChangeAlgorithm creates a new instance with available coin denominations
func NewMoneyChangeAlgorithm(coins []int) *MoneyChangeAlgorithm {
	// Sort coins in descending order for greedy approach
	sortedCoins := make([]int, len(coins))
	copy(sortedCoins, coins)
	sort.Sort(sort.Reverse(sort.IntSlice(sortedCoins)))

	return &MoneyChangeAlgorithm{
		coins: sortedCoins,
	}
}

// ChangeResult represents the result of the change calculation
type ChangeResult struct {
	TotalCoins int
	Breakdown  map[int]int // coin value -> quantity
	Success    bool
	Message    string
}

// CalculateChange finds the optimal combination of coins for a given amount
func (mca *MoneyChangeAlgorithm) CalculateChange(amount int) ChangeResult {
	if amount < 0 {
		return ChangeResult{
			Success: false,
			Message: "Amount cannot be negative",
		}
	}

	if amount == 0 {
		return ChangeResult{
			TotalCoins: 0,
			Breakdown:  make(map[int]int),
			Success:    true,
			Message:    "No change needed",
		}
	}

	remaining := amount
	breakdown := make(map[int]int)

	// Greedy approach: use largest coins first
	for _, coin := range mca.coins {
		if remaining >= coin {
			quantity := remaining / coin
			breakdown[coin] = quantity
			remaining -= quantity * coin
		}
	}

	totalCoins := 0
	for _, quantity := range breakdown {
		totalCoins += quantity
	}

	if remaining > 0 {
		return ChangeResult{
			Success: false,
			Message: fmt.Sprintf("Cannot make exact change. Remaining: %d", remaining),
		}
	}

	return ChangeResult{
		TotalCoins: totalCoins,
		Breakdown:  breakdown,
		Success:    true,
		Message:    fmt.Sprintf("Change calculated with %d coins", totalCoins),
	}
}

// GetAvailableCoins returns the available coin denominations
func (mca *MoneyChangeAlgorithm) GetAvailableCoins() []int {
	return append([]int(nil), mca.coins...)
}
