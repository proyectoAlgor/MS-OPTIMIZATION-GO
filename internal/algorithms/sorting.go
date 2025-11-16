package algorithms

import (
	"sort"
)

// SortingAlgorithm provides various sorting methods for different use cases
type SortingAlgorithm struct{}

// NewSortingAlgorithm creates a new instance
func NewSortingAlgorithm() *SortingAlgorithm {
	return &SortingAlgorithm{}
}

// Product represents a product in the bar system
type Product struct {
	ID       string
	Name     string
	Category string
	Price    float64
	Code     string
}

// QuickSortProducts sorts products using Quick Sort algorithm
func (sa *SortingAlgorithm) QuickSortProducts(products []Product, sortBy string) []Product {
	if len(products) <= 1 {
		return products
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]Product, len(products))
	copy(sorted, products)

	// Use Go's built-in sort with custom comparison
	switch sortBy {
	case "price_asc":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Price < sorted[j].Price
		})
	case "price_desc":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Price > sorted[j].Price
		})
	case "name_asc":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Name < sorted[j].Name
		})
	case "name_desc":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Name > sorted[j].Name
		})
	case "code_asc":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Code < sorted[j].Code
		})
	case "category_asc":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Category < sorted[j].Category
		})
	default:
		// Default to price ascending
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Price < sorted[j].Price
		})
	}

	return sorted
}

// InsertionSortProducts sorts products using Insertion Sort (efficient for small lists)
func (sa *SortingAlgorithm) InsertionSortProducts(products []Product, sortBy string) []Product {
	if len(products) <= 1 {
		return products
	}

	sorted := make([]Product, len(products))
	copy(sorted, products)

	for i := 1; i < len(sorted); i++ {
		key := sorted[i]
		j := i - 1

		// Move elements that are greater than key one position ahead
		for j >= 0 && sa.compareProducts(sorted[j], key, sortBy) > 0 {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}

	return sorted
}

// SelectionSortProducts sorts products using Selection Sort
func (sa *SortingAlgorithm) SelectionSortProducts(products []Product, sortBy string) []Product {
	if len(products) <= 1 {
		return products
	}

	sorted := make([]Product, len(products))
	copy(sorted, products)

	for i := 0; i < len(sorted)-1; i++ {
		minIdx := i
		for j := i + 1; j < len(sorted); j++ {
			if sa.compareProducts(sorted[j], sorted[minIdx], sortBy) < 0 {
				minIdx = j
			}
		}
		sorted[i], sorted[minIdx] = sorted[minIdx], sorted[i]
	}

	return sorted
}

// compareProducts compares two products based on the specified criteria
func (sa *SortingAlgorithm) compareProducts(a, b Product, sortBy string) int {
	switch sortBy {
	case "price_asc":
		if a.Price < b.Price {
			return -1
		} else if a.Price > b.Price {
			return 1
		}
		return 0
	case "price_desc":
		if a.Price > b.Price {
			return -1
		} else if a.Price < b.Price {
			return 1
		}
		return 0
	case "name_asc":
		if a.Name < b.Name {
			return -1
		} else if a.Name > b.Name {
			return 1
		}
		return 0
	case "name_desc":
		if a.Name > b.Name {
			return -1
		} else if a.Name < b.Name {
			return 1
		}
		return 0
	case "code_asc":
		if a.Code < b.Code {
			return -1
		} else if a.Code > b.Code {
			return 1
		}
		return 0
	case "category_asc":
		if a.Category < b.Category {
			return -1
		} else if a.Category > b.Category {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// Table represents a table in the bar system
type Table struct {
	ID       string
	Number   int
	Capacity int
	Status   string
	Location string
}

// SortTables sorts tables using different algorithms
func (sa *SortingAlgorithm) SortTables(tables []Table, sortBy string) []Table {
	if len(tables) <= 1 {
		return tables
	}

	sorted := make([]Table, len(tables))
	copy(sorted, tables)

	switch sortBy {
	case "number_asc":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Number < sorted[j].Number
		})
	case "capacity_asc":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Capacity < sorted[j].Capacity
		})
	case "capacity_desc":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Capacity > sorted[j].Capacity
		})
	case "status":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Status < sorted[j].Status
		})
	default:
		// Default to number ascending
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Number < sorted[j].Number
		})
	}

	return sorted
}
