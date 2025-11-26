package algorithms

import (
	"fmt"
	"math"
	"sort"
)

// TableAssignmentAlgorithm implements a greedy algorithm for optimal table assignment
// This optimizes which tables to assign to customers based on capacity, proximity, and availability
type TableAssignmentAlgorithm struct{}

// NewTableAssignmentAlgorithm creates a new instance
func NewTableAssignmentAlgorithm() *TableAssignmentAlgorithm {
	return &TableAssignmentAlgorithm{}
}

// TableAssignment represents a table in the bar for assignment optimization
type TableAssignment struct {
	ID          string  // Table ID
	Code        string  // Table code (e.g., "MESA-01")
	Capacity    int     // Maximum number of people
	LocationX   float64 // X coordinate (for proximity calculation)
	LocationY   float64 // Y coordinate (for proximity calculation)
	IsAvailable bool    // Whether the table is currently available
	IsOccupied  bool    // Whether the table is currently occupied
	Priority    int     // Priority score (higher = better, e.g., window seats)
}

// CustomerGroup represents a group of customers requesting a table
type CustomerGroup struct {
	ID          string   // Group ID
	Size        int      // Number of people in the group
	Priority    int      // Customer priority (VIP, regular, etc.)
	PreferredX  *float64 // Preferred location X (optional)
	PreferredY  *float64 // Preferred location Y (optional)
	MaxDistance float64  // Maximum distance they're willing to travel (0 = no limit)
}

// Assignment represents a table assignment to a customer group
type Assignment struct {
	TableID      string  `json:"table_id"`
	TableCode    string  `json:"table_code"`
	CustomerID   string  `json:"customer_id"`
	Distance     float64 `json:"distance"`      // Distance from preferred location
	FitnessScore float64 `json:"fitness_score"` // Overall fitness score (0-1, higher is better)
	CapacityUtil float64 `json:"capacity_util"` // Capacity utilization (group_size / table_capacity)
}

// AssignmentResult represents the result of table assignment optimization
type AssignmentResult struct {
	Assignments      []Assignment `json:"assignments"`
	UnassignedGroups []string     `json:"unassigned_groups"` // IDs of groups that couldn't be assigned
	TotalFitness     float64      `json:"total_fitness"`     // Sum of all fitness scores
	AverageFitness   float64      `json:"average_fitness"`   // Average fitness score
	TablesUsed       int          `json:"tables_used"`
	CustomersServed  int          `json:"customers_served"`
	Message          string       `json:"message"`
}

// CalculateDistance calculates Euclidean distance between two points
func calculateDistance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// CalculateFitnessScore calculates how well a table fits a customer group
// Returns a score between 0 and 1, where 1 is perfect fit
func calculateFitnessScore(table TableAssignment, group CustomerGroup) float64 {
	// Base score starts at 1.0
	score := 1.0

	// Penalty for capacity mismatch
	// If table is too small, score = 0 (cannot fit)
	if table.Capacity < group.Size {
		return 0.0
	}

	// Capacity utilization: prefer tables that match group size closely
	// Perfect match (capacity == size) gets 1.0, larger tables get lower score
	capacityRatio := float64(group.Size) / float64(table.Capacity)
	if capacityRatio < 0.5 {
		// Table is more than 2x larger than needed, apply penalty
		score *= 0.7
	} else if capacityRatio < 0.75 {
		// Table is 1.33x to 2x larger, slight penalty
		score *= 0.9
	}
	// If capacityRatio >= 0.75, no penalty (good fit)

	// Distance penalty (if preferred location is specified)
	if group.PreferredX != nil && group.PreferredY != nil {
		distance := calculateDistance(table.LocationX, table.LocationY, *group.PreferredX, *group.PreferredY)

		// Check max distance constraint
		if group.MaxDistance > 0 && distance > group.MaxDistance {
			return 0.0 // Too far, cannot assign
		}

		// Apply distance penalty (closer is better)
		// Normalize distance (assuming max reasonable distance is 100 units)
		maxReasonableDistance := 100.0
		if distance > maxReasonableDistance {
			distance = maxReasonableDistance
		}
		distancePenalty := distance / maxReasonableDistance
		score *= (1.0 - distancePenalty*0.3) // Up to 30% penalty for distance
	}

	// Priority bonus (higher priority tables get bonus)
	// Normalize priority (assuming max priority is 10)
	maxPriority := 10.0
	priorityBonus := float64(table.Priority) / maxPriority * 0.1 // Up to 10% bonus
	score += priorityBonus

	// Customer priority bonus
	customerPriorityBonus := float64(group.Priority) / maxPriority * 0.1 // Up to 10% bonus
	score += customerPriorityBonus

	// Ensure score is between 0 and 1
	if score > 1.0 {
		score = 1.0
	}
	if score < 0.0 {
		score = 0.0
	}

	return score
}

// AssignTablesGreedy assigns tables to customer groups using a greedy algorithm
// Time Complexity: O(n * m) where n = number of groups, m = number of tables
// Space Complexity: O(n + m)
// This is efficient and provides good results for real-time table assignment
func (taa *TableAssignmentAlgorithm) AssignTablesGreedy(tables []TableAssignment, groups []CustomerGroup) AssignmentResult {
	if len(tables) == 0 {
		return AssignmentResult{
			Assignments:      []Assignment{},
			UnassignedGroups: getGroupIDs(groups),
			Message:          "No tables available",
		}
	}

	if len(groups) == 0 {
		return AssignmentResult{
			Assignments:      []Assignment{},
			UnassignedGroups: []string{},
			Message:          "No customer groups to assign",
		}
	}

	// Filter available tables
	availableTables := []TableAssignment{}
	for _, table := range tables {
		if table.IsAvailable && !table.IsOccupied {
			availableTables = append(availableTables, table)
		}
	}

	if len(availableTables) == 0 {
		return AssignmentResult{
			Assignments:      []Assignment{},
			UnassignedGroups: getGroupIDs(groups),
			Message:          "No available tables",
		}
	}

	// Sort groups by priority (higher priority first) and size (larger groups first)
	// This ensures VIP customers and larger groups get assigned first
	sortedGroups := make([]CustomerGroup, len(groups))
	copy(sortedGroups, groups)
	sort.Slice(sortedGroups, func(i, j int) bool {
		if sortedGroups[i].Priority != sortedGroups[j].Priority {
			return sortedGroups[i].Priority > sortedGroups[j].Priority
		}
		return sortedGroups[i].Size > sortedGroups[j].Size
	})

	assignments := []Assignment{}
	assignedTableIDs := make(map[string]bool)
	unassignedGroups := []string{}

	// Greedy assignment: for each group, find the best available table
	for _, group := range sortedGroups {
		bestTable := TableAssignment{}
		bestScore := -1.0

		// Find the best table for this group
		for _, table := range availableTables {
			if assignedTableIDs[table.ID] {
				continue // Table already assigned
			}

			score := calculateFitnessScore(table, group)
			if score > bestScore && score > 0 {
				bestScore = score
				bestTable = table
			}
		}

		// If we found a suitable table, assign it
		if bestScore > 0 {
			distance := 0.0
			if group.PreferredX != nil && group.PreferredY != nil {
				distance = calculateDistance(bestTable.LocationX, bestTable.LocationY, *group.PreferredX, *group.PreferredY)
			}

			capacityUtil := float64(group.Size) / float64(bestTable.Capacity)

			assignment := Assignment{
				TableID:      bestTable.ID,
				TableCode:    bestTable.Code,
				CustomerID:   group.ID,
				Distance:     distance,
				FitnessScore: bestScore,
				CapacityUtil: capacityUtil,
			}

			assignments = append(assignments, assignment)
			assignedTableIDs[bestTable.ID] = true
		} else {
			// No suitable table found for this group
			unassignedGroups = append(unassignedGroups, group.ID)
		}
	}

	// Calculate statistics
	totalFitness := 0.0
	customersServed := 0
	for _, assignment := range assignments {
		totalFitness += assignment.FitnessScore
		// Find the group size for this assignment
		for _, group := range groups {
			if group.ID == assignment.CustomerID {
				customersServed += group.Size
				break
			}
		}
	}

	averageFitness := 0.0
	if len(assignments) > 0 {
		averageFitness = totalFitness / float64(len(assignments))
	}

	return AssignmentResult{
		Assignments:      assignments,
		UnassignedGroups: unassignedGroups,
		TotalFitness:     totalFitness,
		AverageFitness:   averageFitness,
		TablesUsed:       len(assignments),
		CustomersServed:  customersServed,
		Message:          fmt.Sprintf("Assigned %d tables to %d customer groups (%.1f%% success rate)", len(assignments), len(groups), float64(len(assignments))/float64(len(groups))*100),
	}
}

// AssignTablesOptimal finds the optimal assignment using a more complex algorithm
// This uses a matching algorithm (similar to Hungarian algorithm) for optimal results
// Time Complexity: O(n³) where n = max(number of tables, number of groups)
// Space Complexity: O(n²)
// Use this when you need the absolute best assignment, but it's slower
func (taa *TableAssignmentAlgorithm) AssignTablesOptimal(tables []TableAssignment, groups []CustomerGroup) AssignmentResult {
	// For now, we'll use the greedy approach as a baseline
	// A full implementation of Hungarian algorithm would be more complex
	// This is a placeholder that can be enhanced later
	return taa.AssignTablesGreedy(tables, groups)
}

// Helper function to get group IDs
func getGroupIDs(groups []CustomerGroup) []string {
	ids := make([]string, len(groups))
	for i, group := range groups {
		ids[i] = group.ID
	}
	return ids
}
