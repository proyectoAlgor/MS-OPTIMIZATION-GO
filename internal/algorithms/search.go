package algorithms

import (
	"fmt"
	"strings"
)

// SearchAlgorithm provides various search methods for different use cases
type SearchAlgorithm struct{}

// NewSearchAlgorithm creates a new instance
func NewSearchAlgorithm() *SearchAlgorithm {
	return &SearchAlgorithm{}
}

// BinarySearchResult represents the result of a binary search
type BinarySearchResult struct {
	Found   bool
	Index   int
	Message string
}

// BinarySearchProducts searches for a product by price using binary search
func (sa *SearchAlgorithm) BinarySearchProducts(products []Product, targetPrice float64) BinarySearchResult {
	if len(products) == 0 {
		return BinarySearchResult{
			Found:   false,
			Index:   -1,
			Message: "Empty product list",
		}
	}

	// First, we need to sort the products by price for binary search
	sortingAlgo := NewSortingAlgorithm()
	sortedProducts := sortingAlgo.QuickSortProducts(products, "price_asc")

	left, right := 0, len(sortedProducts)-1

	for left <= right {
		mid := left + (right-left)/2

		if sortedProducts[mid].Price == targetPrice {
			return BinarySearchResult{
				Found: true,
				Index: mid,
				Message: fmt.Sprintf("Found product '%s' at index %d with price %.2f",
					sortedProducts[mid].Name, mid, targetPrice),
			}
		}

		if sortedProducts[mid].Price < targetPrice {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return BinarySearchResult{
		Found:   false,
		Index:   -1,
		Message: fmt.Sprintf("No product found with price %.2f", targetPrice),
	}
}

// BinarySearchProductsByPriceRange finds products within a price range
func (sa *SearchAlgorithm) BinarySearchProductsByPriceRange(products []Product, minPrice, maxPrice float64) []Product {
	if len(products) == 0 {
		return []Product{}
	}

	// Sort products by price
	sortingAlgo := NewSortingAlgorithm()
	sortedProducts := sortingAlgo.QuickSortProducts(products, "price_asc")

	var result []Product
	for _, product := range sortedProducts {
		if product.Price >= minPrice && product.Price <= maxPrice {
			result = append(result, product)
		}
	}

	return result
}

// ReverseString reverses a string (useful for searching in reverse order)
func (sa *SearchAlgorithm) ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// SearchProductsByName searches for products by name (case-insensitive)
func (sa *SearchAlgorithm) SearchProductsByName(products []Product, searchTerm string) []Product {
	if searchTerm == "" {
		return products
	}

	searchTerm = strings.ToLower(searchTerm)
	var result []Product

	for _, product := range products {
		productName := strings.ToLower(product.Name)
		if strings.Contains(productName, searchTerm) {
			result = append(result, product)
		}
	}

	return result
}

// SearchProductsByCode searches for products by exact code match
func (sa *SearchAlgorithm) SearchProductsByCode(products []Product, code string) *Product {
	for _, product := range products {
		if strings.EqualFold(product.Code, code) {
			return &product
		}
	}
	return nil
}

// SumProductPrices calculates the total price of a list of products
func (sa *SearchAlgorithm) SumProductPrices(products []Product) float64 {
	total := 0.0
	for _, product := range products {
		total += product.Price
	}
	return total
}

// SumProductPricesRecursive calculates the total price recursively
func (sa *SearchAlgorithm) SumProductPricesRecursive(products []Product) float64 {
	if len(products) == 0 {
		return 0.0
	}
	if len(products) == 1 {
		return products[0].Price
	}

	// Divide and conquer approach
	mid := len(products) / 2
	leftSum := sa.SumProductPricesRecursive(products[:mid])
	rightSum := sa.SumProductPricesRecursive(products[mid:])

	return leftSum + rightSum
}

// Order represents an order in the bar system
type Order struct {
	ID       string
	TableID  string
	Products []Product
	Total    float64
	Status   string
}

// CalculateOrderTotal calculates the total of an order
func (sa *SearchAlgorithm) CalculateOrderTotal(order Order) float64 {
	total := 0.0
	for _, product := range order.Products {
		total += product.Price
	}
	return total
}

// FindMostExpensiveProduct finds the most expensive product in a list
func (sa *SearchAlgorithm) FindMostExpensiveProduct(products []Product) *Product {
	if len(products) == 0 {
		return nil
	}

	mostExpensive := &products[0]
	for i := 1; i < len(products); i++ {
		if products[i].Price > mostExpensive.Price {
			mostExpensive = &products[i]
		}
	}

	return mostExpensive
}

// FindCheapestProduct finds the cheapest product in a list
func (sa *SearchAlgorithm) FindCheapestProduct(products []Product) *Product {
	if len(products) == 0 {
		return nil
	}

	cheapest := &products[0]
	for i := 1; i < len(products); i++ {
		if products[i].Price < cheapest.Price {
			cheapest = &products[i]
		}
	}

	return cheapest
}

