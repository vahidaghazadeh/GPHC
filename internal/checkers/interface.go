package checkers

import (
	"github.com/opsource/gphc/pkg/types"
)

// Checker defines the interface that all checkers must implement
type Checker interface {
	// Name returns the name of the checker
	Name() string
	
	// ID returns the unique identifier of the checker
	ID() string
	
	// Category returns the category this checker belongs to
	Category() types.Category
	
	// Check performs the actual check and returns the result
	Check(data *types.RepositoryData) *types.CheckResult
	
	// Weight returns the weight/importance of this checker (1-10)
	Weight() int
}

// BaseChecker provides common functionality for all checkers
type BaseChecker struct {
	name     string
	id       string
	category types.Category
	weight   int
}

// NewBaseChecker creates a new base checker
func NewBaseChecker(name, id string, category types.Category, weight int) BaseChecker {
	return BaseChecker{
		name:     name,
		id:       id,
		category: category,
		weight:   weight,
	}
}

// Name returns the checker name
func (b BaseChecker) Name() string {
	return b.name
}

// ID returns the checker ID
func (b BaseChecker) ID() string {
	return b.id
}

// Category returns the checker category
func (b BaseChecker) Category() types.Category {
	return b.category
}

// Weight returns the checker weight
func (b BaseChecker) Weight() int {
	return b.weight
}
