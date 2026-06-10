package scene

import (
	"testing"
)

func TestSortScenes(t *testing.T) {
	tests := []struct {
		name     string
		scenes   map[string]*Scene
		expected []struct {
			chapter     int
			sceneNumber int
			name        string
		}
	}{
		{
			name:   "empty map",
			scenes: make(map[string]*Scene),
			expected: []struct {
				chapter     int
				sceneNumber int
				name        string
			}{},
		},
		{
			name: "single scene",
			scenes: map[string]*Scene{
				"scene1": {
					ID:          "scene1",
					Name:        "Opening",
					Chapter:     1,
					SceneNumber: 1,
				},
			},
			expected: []struct {
				chapter     int
				sceneNumber int
				name        string
			}{
				{1, 1, "Opening"},
			},
		},
		{
			name: "multiple scenes in order",
			scenes: map[string]*Scene{
				"scene1": {
					ID:          "scene1",
					Name:        "Scene 1",
					Chapter:     1,
					SceneNumber: 1,
				},
				"scene2": {
					ID:          "scene2",
					Name:        "Scene 2",
					Chapter:     1,
					SceneNumber: 2,
				},
				"scene3": {
					ID:          "scene3",
					Name:        "Scene 3",
					Chapter:     2,
					SceneNumber: 1,
				},
			},
			expected: []struct {
				chapter     int
				sceneNumber int
				name        string
			}{
				{1, 1, "Scene 1"},
				{1, 2, "Scene 2"},
				{2, 1, "Scene 3"},
			},
		},
		{
			name: "multiple scenes out of order",
			scenes: map[string]*Scene{
				"scene3": {
					ID:          "scene3",
					Name:        "Scene 3",
					Chapter:     2,
					SceneNumber: 1,
				},
				"scene1": {
					ID:          "scene1",
					Name:        "Scene 1",
					Chapter:     1,
					SceneNumber: 1,
				},
				"scene2": {
					ID:          "scene2",
					Name:        "Scene 2",
					Chapter:     1,
					SceneNumber: 2,
				},
			},
			expected: []struct {
				chapter     int
				sceneNumber int
				name        string
			}{
				{1, 1, "Scene 1"},
				{1, 2, "Scene 2"},
				{2, 1, "Scene 3"},
			},
		},
		{
			name: "complex ordering",
			scenes: map[string]*Scene{
				"scene5": {
					ID:          "scene5",
					Name:        "Scene 5",
					Chapter:     3,
					SceneNumber: 1,
				},
				"scene3": {
					ID:          "scene3",
					Name:        "Scene 3",
					Chapter:     1,
					SceneNumber: 3,
				},
				"scene1": {
					ID:          "scene1",
					Name:        "Scene 1",
					Chapter:     1,
					SceneNumber: 1,
				},
				"scene4": {
					ID:          "scene4",
					Name:        "Scene 4",
					Chapter:     2,
					SceneNumber: 1,
				},
				"scene2": {
					ID:          "scene2",
					Name:        "Scene 2",
					Chapter:     1,
					SceneNumber: 2,
				},
			},
			expected: []struct {
				chapter     int
				sceneNumber int
				name        string
			}{
				{1, 1, "Scene 1"},
				{1, 2, "Scene 2"},
				{1, 3, "Scene 3"},
				{2, 1, "Scene 4"},
				{3, 1, "Scene 5"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortScenes(tt.scenes)

			if len(result) != len(tt.expected) {
				t.Errorf("SortScenes() returned %d scenes, expected %d", len(result), len(tt.expected))
				return
			}

			for i, exp := range tt.expected {
				if result[i].Chapter != exp.chapter {
					t.Errorf("Scene %d: got chapter %d, expected %d", i, result[i].Chapter, exp.chapter)
				}
				if result[i].SceneNumber != exp.sceneNumber {
					t.Errorf("Scene %d: got scene number %d, expected %d", i, result[i].SceneNumber, exp.sceneNumber)
				}
				if result[i].Name != exp.name {
					t.Errorf("Scene %d: got name %q, expected %q", i, result[i].Name, exp.name)
				}
			}
		})
	}
}
