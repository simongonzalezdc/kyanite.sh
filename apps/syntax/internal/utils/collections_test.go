package utils

import (
	"testing"

	"github.com/kyanite/syntax/internal/character"
	"github.com/kyanite/syntax/internal/location"
	"github.com/kyanite/syntax/internal/scene"
)

func TestGetSceneAtIndex(t *testing.T) {
	scenes := map[string]*scene.Scene{
		"scene1": {ID: "scene1", Name: "Scene 1", Chapter: 1, SceneNumber: 1},
		"scene2": {ID: "scene2", Name: "Scene 2", Chapter: 1, SceneNumber: 2},
		"scene3": {ID: "scene3", Name: "Scene 3", Chapter: 2, SceneNumber: 1},
	}

	tests := []struct {
		name     string
		index    int
		expected string // Expected scene name, or empty if nil expected
	}{
		{"first scene", 0, "Scene 1"},
		{"second scene", 1, "Scene 2"},
		{"third scene", 2, "Scene 3"},
		{"index too high", 3, ""},
		{"negative index", -1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSceneAtIndex(scenes, tt.index)
			if tt.expected == "" {
				if result != nil {
					t.Errorf("GetSceneAtIndex(%d) = %v, expected nil", tt.index, result)
				}
			} else {
				if result == nil {
					t.Errorf("GetSceneAtIndex(%d) = nil, expected scene with name %q", tt.index, tt.expected)
				} else if result.Name != tt.expected {
					t.Errorf("GetSceneAtIndex(%d) = %q, expected %q", tt.index, result.Name, tt.expected)
				}
			}
		})
	}
}

func TestGetCharacterAtIndex(t *testing.T) {
	characters := map[string]*character.Character{
		"char1": {ID: "char1", Name: "Alice"},
		"char2": {ID: "char2", Name: "Bob"},
	}

	tests := []struct {
		name     string
		index    int
		expected string
	}{
		{"first character", 0, "Alice"},
		{"second character", 1, "Bob"},
		{"index too high", 2, ""},
		{"negative index", -1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCharacterAtIndex(characters, tt.index)
			if tt.expected == "" {
				if result != nil {
					t.Errorf("GetCharacterAtIndex(%d) = %v, expected nil", tt.index, result)
				}
			} else {
				if result == nil {
					t.Errorf("GetCharacterAtIndex(%d) = nil, expected character with name %q", tt.index, tt.expected)
				} else if result.Name != tt.expected {
					t.Errorf("GetCharacterAtIndex(%d) = %q, expected %q", tt.index, result.Name, tt.expected)
				}
			}
		})
	}
}

func TestGetLocationAtIndex(t *testing.T) {
	locations := map[string]*location.Location{
		"loc1": {ID: "loc1", Name: "Forest"},
		"loc2": {ID: "loc2", Name: "Castle"},
	}

	// Test valid indices - we can't predict order but can verify we get valid results
	t.Run("first location", func(t *testing.T) {
		result := GetLocationAtIndex(locations, 0)
		if result == nil {
			t.Error("GetLocationAtIndex(0) = nil, expected a location")
		} else if result.Name != "Forest" && result.Name != "Castle" {
			t.Errorf("GetLocationAtIndex(0) = %q, expected Forest or Castle", result.Name)
		}
	})

	t.Run("second location", func(t *testing.T) {
		result := GetLocationAtIndex(locations, 1)
		if result == nil {
			t.Error("GetLocationAtIndex(1) = nil, expected a location")
		} else if result.Name != "Forest" && result.Name != "Castle" {
			t.Errorf("GetLocationAtIndex(1) = %q, expected Forest or Castle", result.Name)
		}
	})

	t.Run("index too high", func(t *testing.T) {
		result := GetLocationAtIndex(locations, 2)
		if result != nil {
			t.Errorf("GetLocationAtIndex(2) = %v, expected nil", result)
		}
	})

	t.Run("negative index", func(t *testing.T) {
		result := GetLocationAtIndex(locations, -1)
		if result != nil {
			t.Errorf("GetLocationAtIndex(-1) = %v, expected nil", result)
		}
	})
}

func TestFindSceneIDAtIndex(t *testing.T) {
	scenes := map[string]*scene.Scene{
		"scene1": {ID: "scene1", Name: "Scene 1", Chapter: 1, SceneNumber: 1},
		"scene2": {ID: "scene2", Name: "Scene 2", Chapter: 1, SceneNumber: 2},
	}

	tests := []struct {
		name     string
		index    int
		expected string
	}{
		{"first scene", 0, "scene1"},
		{"second scene", 1, "scene2"},
		{"index too high", 2, ""},
		{"negative index", -1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindSceneIDAtIndex(scenes, tt.index)
			if result != tt.expected {
				t.Errorf("FindSceneIDAtIndex(%d) = %q, expected %q", tt.index, result, tt.expected)
			}
		})
	}
}

func TestFindCharacterIDAtIndex(t *testing.T) {
	characters := map[string]*character.Character{
		"char1": {ID: "char1", Name: "Alice"},
		"char2": {ID: "char2", Name: "Bob"},
	}

	// Test valid indices - we can't predict order but can verify we get valid IDs
	t.Run("first character", func(t *testing.T) {
		result := FindCharacterIDAtIndex(characters, 0)
		if result != "char1" && result != "char2" {
			t.Errorf("FindCharacterIDAtIndex(0) = %q, expected char1 or char2", result)
		}
	})

	t.Run("second character", func(t *testing.T) {
		result := FindCharacterIDAtIndex(characters, 1)
		if result != "char1" && result != "char2" {
			t.Errorf("FindCharacterIDAtIndex(1) = %q, expected char1 or char2", result)
		}
	})

	t.Run("index too high", func(t *testing.T) {
		result := FindCharacterIDAtIndex(characters, 2)
		if result != "" {
			t.Errorf("FindCharacterIDAtIndex(2) = %q, expected empty string", result)
		}
	})

	t.Run("negative index", func(t *testing.T) {
		result := FindCharacterIDAtIndex(characters, -1)
		if result != "" {
			t.Errorf("FindCharacterIDAtIndex(-1) = %q, expected empty string", result)
		}
	})
}

func TestFindLocationIDAtIndex(t *testing.T) {
	locations := map[string]*location.Location{
		"loc1": {ID: "loc1", Name: "Forest"},
		"loc2": {ID: "loc2", Name: "Castle"},
	}

	// Test valid indices - we can't predict order but can verify we get valid IDs
	t.Run("first location", func(t *testing.T) {
		result := FindLocationIDAtIndex(locations, 0)
		if result != "loc1" && result != "loc2" {
			t.Errorf("FindLocationIDAtIndex(0) = %q, expected loc1 or loc2", result)
		}
	})

	t.Run("second location", func(t *testing.T) {
		result := FindLocationIDAtIndex(locations, 1)
		if result != "loc1" && result != "loc2" {
			t.Errorf("FindLocationIDAtIndex(1) = %q, expected loc1 or loc2", result)
		}
	})

	t.Run("index too high", func(t *testing.T) {
		result := FindLocationIDAtIndex(locations, 2)
		if result != "" {
			t.Errorf("FindLocationIDAtIndex(2) = %q, expected empty string", result)
		}
	})

	t.Run("negative index", func(t *testing.T) {
		result := FindLocationIDAtIndex(locations, -1)
		if result != "" {
			t.Errorf("FindLocationIDAtIndex(-1) = %q, expected empty string", result)
		}
	})
}

func TestEmptyMaps(t *testing.T) {
	t.Run("GetSceneAtIndex with empty map", func(t *testing.T) {
		result := GetSceneAtIndex(make(map[string]*scene.Scene), 0)
		if result != nil {
			t.Errorf("Expected nil for empty map, got %v", result)
		}
	})

	t.Run("GetCharacterAtIndex with empty map", func(t *testing.T) {
		result := GetCharacterAtIndex(make(map[string]*character.Character), 0)
		if result != nil {
			t.Errorf("Expected nil for empty map, got %v", result)
		}
	})

	t.Run("GetLocationAtIndex with empty map", func(t *testing.T) {
		result := GetLocationAtIndex(make(map[string]*location.Location), 0)
		if result != nil {
			t.Errorf("Expected nil for empty map, got %v", result)
		}
	})

	t.Run("FindSceneIDAtIndex with empty map", func(t *testing.T) {
		result := FindSceneIDAtIndex(make(map[string]*scene.Scene), 0)
		if result != "" {
			t.Errorf("Expected empty string for empty map, got %q", result)
		}
	})

	t.Run("FindCharacterIDAtIndex with empty map", func(t *testing.T) {
		result := FindCharacterIDAtIndex(make(map[string]*character.Character), 0)
		if result != "" {
			t.Errorf("Expected empty string for empty map, got %q", result)
		}
	})

	t.Run("FindLocationIDAtIndex with empty map", func(t *testing.T) {
		result := FindLocationIDAtIndex(make(map[string]*location.Location), 0)
		if result != "" {
			t.Errorf("Expected empty string for empty map, got %q", result)
		}
	})
}
