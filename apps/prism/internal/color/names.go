package color

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/kyanite/prism/internal/data"
)

// NamedColor represents a named color
type NamedColor struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

// ColorDatabase holds all named colors
type ColorDatabase struct {
	Colors []NamedColor `json:"colors"`
}

var (
	db *ColorDatabase
	mu sync.RWMutex
)

// LoadNamedColors loads the named colors database
func LoadNamedColors() error {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		return nil // Already loaded
	}

	db = &ColorDatabase{}
	err := json.Unmarshal(data.ColorsJSON, db)
	if err != nil {
		return err
	}

	return nil
}

// SearchColors searches for colors by name (fuzzy search)
func SearchColors(query string) []NamedColor {
	mu.RLock()
	if db == nil {
		mu.RUnlock()
		if err := LoadNamedColors(); err != nil {
			return []NamedColor{}
		}
		mu.RLock()
	}
	defer mu.RUnlock()

	query = strings.ToLower(query)
	// Pre-allocate with reasonable capacity for search results
	results := make([]NamedColor, 0, 20)

	// Exact matches first
	for _, c := range db.Colors {
		if strings.ToLower(c.Name) == query {
			results = append(results, c)
		}
	}

	// Prefix matches
	for _, c := range db.Colors {
		if strings.HasPrefix(strings.ToLower(c.Name), query) && !contains(results, c) {
			results = append(results, c)
		}
	}

	// Contains matches
	for _, c := range db.Colors {
		if strings.Contains(strings.ToLower(c.Name), query) && !contains(results, c) {
			results = append(results, c)
		}
	}

	// Fuzzy matches (1-2 character differences)
	for _, c := range db.Colors {
		if fuzzyMatch(strings.ToLower(c.Name), query) && !contains(results, c) {
			results = append(results, c)
		}
	}

	// Limit results
	if len(results) > 10 {
		results = results[:10]
	}

	return results
}

// GetColorByName returns a color by exact name
func GetColorByName(name string) (*NamedColor, error) {
	mu.RLock()
	if db == nil {
		mu.RUnlock()
		if err := LoadNamedColors(); err != nil {
			return nil, err
		}
		mu.RLock()
	}
	defer mu.RUnlock()

	name = strings.ToLower(name)
	for _, c := range db.Colors {
		if strings.ToLower(c.Name) == name {
			return &c, nil
		}
	}

	return nil, nil
}

// AllNamedColors returns all named colors
func AllNamedColors() []NamedColor {
	mu.RLock()
	if db == nil {
		mu.RUnlock()
		if err := LoadNamedColors(); err != nil {
			return []NamedColor{}
		}
		mu.RLock()
	}
	defer mu.RUnlock()

	return db.Colors
}

// contains checks if a color is in the results
func contains(results []NamedColor, color NamedColor) bool {
	for _, r := range results {
		if r.Name == color.Name {
			return true
		}
	}
	return false
}

// fuzzyMatch performs simple fuzzy matching (Levenshtein distance <= 2)
func fuzzyMatch(str1, str2 string) bool {
	if len(str1) == 0 || len(str2) == 0 {
		return false
	}

	// Simple length check
	if abs(len(str1)-len(str2)) > 2 {
		return false
	}

	distance := levenshtein(str1, str2)
	return distance <= 2 && distance > 0
}

// levenshtein calculates Levenshtein distance
func levenshtein(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}

	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min3(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
