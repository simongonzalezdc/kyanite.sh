package utils

import (
	"github.com/kyanite/syntax/internal/character"
	"github.com/kyanite/syntax/internal/location"
	"github.com/kyanite/syntax/internal/scene"
)

// GetSceneAtIndex returns the scene at the given index after sorting
// Returns nil if index is out of bounds
func GetSceneAtIndex(scenes map[string]*scene.Scene, index int) *scene.Scene {
	sorted := scene.SortScenes(scenes)
	if index < 0 || index >= len(sorted) {
		return nil
	}
	return sorted[index]
}

// GetCharacterAtIndex returns the character at the given index
// Note: Characters don't have a natural sort order, so order may vary
// Returns nil if index is out of bounds
func GetCharacterAtIndex(characters map[string]*character.Character, index int) *character.Character {
	if index < 0 || index >= len(characters) {
		return nil
	}

	// Convert map to slice
	charList := make([]*character.Character, 0, len(characters))
	for _, char := range characters {
		charList = append(charList, char)
	}

	return charList[index]
}

// GetLocationAtIndex returns the location at the given index
// Note: Locations don't have a natural sort order, so order may vary
// Returns nil if index is out of bounds
func GetLocationAtIndex(locations map[string]*location.Location, index int) *location.Location {
	if index < 0 || index >= len(locations) {
		return nil
	}

	// Convert map to slice
	locList := make([]*location.Location, 0, len(locations))
	for _, loc := range locations {
		locList = append(locList, loc)
	}

	return locList[index]
}

// FindSceneByID returns the scene ID for a scene at the given index after sorting
// Returns empty string if index is out of bounds
func FindSceneIDAtIndex(scenes map[string]*scene.Scene, index int) string {
	sc := GetSceneAtIndex(scenes, index)
	if sc == nil {
		return ""
	}
	return sc.ID
}

// FindCharacterIDAtIndex returns the character ID at the given index
// Returns empty string if index is out of bounds
func FindCharacterIDAtIndex(characters map[string]*character.Character, index int) string {
	char := GetCharacterAtIndex(characters, index)
	if char == nil {
		return ""
	}
	return char.ID
}

// FindLocationIDAtIndex returns the location ID at the given index
// Returns empty string if index is out of bounds
func FindLocationIDAtIndex(locations map[string]*location.Location, index int) string {
	loc := GetLocationAtIndex(locations, index)
	if loc == nil {
		return ""
	}
	return loc.ID
}
