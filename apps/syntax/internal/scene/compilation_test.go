package scene

import (
	"strings"
	"testing"
	"time"
)

func TestCompileScenes_Empty(t *testing.T) {
	scenes := make(map[string]*Scene)
	report := CompileScenes(scenes)

	if report.TotalScenes != 0 {
		t.Errorf("TotalScenes = %d, expected 0", report.TotalScenes)
	}

	if len(report.Issues) != 1 {
		t.Errorf("Issues count = %d, expected 1", len(report.Issues))
	}

	if len(report.Issues) > 0 {
		issue := report.Issues[0]
		if issue.Type != "empty" {
			t.Errorf("Issue type = %q, expected 'empty'", issue.Type)
		}
		if issue.Severity != "info" {
			t.Errorf("Issue severity = %q, expected 'info'", issue.Severity)
		}
	}

	if report.ChapterCount != 0 {
		t.Errorf("ChapterCount = %d, expected 0", report.ChapterCount)
	}
}

func TestCompileScenes_Valid(t *testing.T) {
	scenes := map[string]*Scene{
		"scene1": {
			ID:          "scene1",
			Name:        "Opening",
			Chapter:     1,
			SceneNumber: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		"scene2": {
			ID:          "scene2",
			Name:        "Conflict",
			Chapter:     1,
			SceneNumber: 2,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		"scene3": {
			ID:          "scene3",
			Name:        "Resolution",
			Chapter:     2,
			SceneNumber: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	report := CompileScenes(scenes)

	if report.TotalScenes != 3 {
		t.Errorf("TotalScenes = %d, expected 3", report.TotalScenes)
	}

	if len(report.Issues) != 0 {
		t.Errorf("Issues count = %d, expected 0 for valid scenes", len(report.Issues))
	}

	if report.ChapterCount != 2 {
		t.Errorf("ChapterCount = %d, expected 2", report.ChapterCount)
	}

	if len(report.Orphans) != 0 {
		t.Errorf("Orphans count = %d, expected 0", len(report.Orphans))
	}

	if len(report.Gaps) != 0 {
		t.Errorf("Gaps count = %d, expected 0", len(report.Gaps))
	}

	if len(report.Duplicates) != 0 {
		t.Errorf("Duplicates count = %d, expected 0", len(report.Duplicates))
	}
}

func TestCompileScenes_Orphans(t *testing.T) {
	scenes := map[string]*Scene{
		"scene1": {
			ID:          "scene1",
			Name:        "Valid Scene",
			Chapter:     1,
			SceneNumber: 1,
		},
		"orphan1": {
			ID:          "orphan1",
			Name:        "Orphan No Chapter",
			Chapter:     0,
			SceneNumber: 1,
		},
		"orphan2": {
			ID:          "orphan2",
			Name:        "Orphan No Scene Number",
			Chapter:     1,
			SceneNumber: 0,
		},
	}

	report := CompileScenes(scenes)

	if len(report.Orphans) != 2 {
		t.Errorf("Orphans count = %d, expected 2", len(report.Orphans))
	}

	orphanIssues := 0
	for _, issue := range report.Issues {
		if issue.Type == "orphan" {
			orphanIssues++
			if issue.Severity != "warning" {
				t.Errorf("Orphan issue severity = %q, expected 'warning'", issue.Severity)
			}
		}
	}

	if orphanIssues != 2 {
		t.Errorf("Orphan issues = %d, expected 2", orphanIssues)
	}
}

func TestCompileScenes_Gaps(t *testing.T) {
	scenes := map[string]*Scene{
		"scene1": {
			ID:          "scene1",
			Name:        "Scene 1",
			Chapter:     1,
			SceneNumber: 1,
		},
		"scene3": {
			ID:          "scene3",
			Name:        "Scene 3",
			Chapter:     1,
			SceneNumber: 3,
		},
		"scene5": {
			ID:          "scene5",
			Name:        "Scene 5",
			Chapter:     1,
			SceneNumber: 5,
		},
	}

	report := CompileScenes(scenes)

	// Should detect gaps at scene 2 and scene 4
	if len(report.Gaps) != 2 {
		t.Errorf("Gaps count = %d, expected 2", len(report.Gaps))
	}

	gapIssues := 0
	for _, issue := range report.Issues {
		if issue.Type == "gap" {
			gapIssues++
			if issue.Severity != "warning" {
				t.Errorf("Gap issue severity = %q, expected 'warning'", issue.Severity)
			}
		}
	}

	if gapIssues != 2 {
		t.Errorf("Gap issues = %d, expected 2", gapIssues)
	}

	// Check specific gap information
	foundGap2 := false
	foundGap4 := false
	for _, gap := range report.Gaps {
		if gap.Chapter == 1 && gap.SceneNumber == 2 {
			foundGap2 = true
		}
		if gap.Chapter == 1 && gap.SceneNumber == 4 {
			foundGap4 = true
		}
	}

	if !foundGap2 {
		t.Error("Did not find gap for scene 2")
	}
	if !foundGap4 {
		t.Error("Did not find gap for scene 4")
	}
}

func TestCompileScenes_GapAtStart(t *testing.T) {
	scenes := map[string]*Scene{
		"scene3": {
			ID:          "scene3",
			Name:        "Scene 3",
			Chapter:     1,
			SceneNumber: 3,
		},
	}

	report := CompileScenes(scenes)

	// Should detect gaps at scene 1 and scene 2
	if len(report.Gaps) != 2 {
		t.Errorf("Gaps count = %d, expected 2 (scenes 1 and 2 missing)", len(report.Gaps))
	}

	foundGap1 := false
	foundGap2 := false
	for _, gap := range report.Gaps {
		if gap.Chapter == 1 && gap.SceneNumber == 1 {
			foundGap1 = true
			if !strings.Contains(gap.Message, "at start") {
				t.Error("Gap at scene 1 should mention 'at start'")
			}
		}
		if gap.Chapter == 1 && gap.SceneNumber == 2 {
			foundGap2 = true
		}
	}

	if !foundGap1 {
		t.Error("Did not find gap for scene 1")
	}
	if !foundGap2 {
		t.Error("Did not find gap for scene 2")
	}
}

func TestCompileScenes_Duplicates(t *testing.T) {
	scenes := map[string]*Scene{
		"scene1a": {
			ID:          "scene1a",
			Name:        "Scene 1 Version A",
			Chapter:     1,
			SceneNumber: 1,
		},
		"scene1b": {
			ID:          "scene1b",
			Name:        "Scene 1 Version B",
			Chapter:     1,
			SceneNumber: 1,
		},
		"scene2a": {
			ID:          "scene2a",
			Name:        "Scene 2 Version A",
			Chapter:     1,
			SceneNumber: 2,
		},
		"scene2b": {
			ID:          "scene2b",
			Name:        "Scene 2 Version B",
			Chapter:     1,
			SceneNumber: 2,
		},
		"scene2c": {
			ID:          "scene2c",
			Name:        "Scene 2 Version C",
			Chapter:     1,
			SceneNumber: 2,
		},
	}

	report := CompileScenes(scenes)

	// Should detect 2 duplicate groups
	if len(report.Duplicates) != 2 {
		t.Errorf("Duplicates count = %d, expected 2", len(report.Duplicates))
	}

	duplicateIssues := 0
	for _, issue := range report.Issues {
		if issue.Type == "duplicate" {
			duplicateIssues++
			if issue.Severity != "error" {
				t.Errorf("Duplicate issue severity = %q, expected 'error'", issue.Severity)
			}
		}
	}

	if duplicateIssues != 2 {
		t.Errorf("Duplicate issues = %d, expected 2", duplicateIssues)
	}

	// Check that scene 1 has 2 duplicates and scene 2 has 3
	for _, dup := range report.Duplicates {
		if dup.SceneNumber == 1 && len(dup.SceneIDs) != 2 {
			t.Errorf("Scene 1 duplicates = %d, expected 2", len(dup.SceneIDs))
		}
		if dup.SceneNumber == 2 && len(dup.SceneIDs) != 3 {
			t.Errorf("Scene 2 duplicates = %d, expected 3", len(dup.SceneIDs))
		}
	}
}

func TestCompileScenes_MultipleChapters(t *testing.T) {
	scenes := map[string]*Scene{
		"c1s1": {ID: "c1s1", Chapter: 1, SceneNumber: 1},
		"c1s2": {ID: "c1s2", Chapter: 1, SceneNumber: 2},
		"c2s1": {ID: "c2s1", Chapter: 2, SceneNumber: 1},
		"c3s1": {ID: "c3s1", Chapter: 3, SceneNumber: 1},
		"c3s2": {ID: "c3s2", Chapter: 3, SceneNumber: 2},
	}

	report := CompileScenes(scenes)

	if report.ChapterCount != 3 {
		t.Errorf("ChapterCount = %d, expected 3", report.ChapterCount)
	}

	if len(report.Issues) != 0 {
		t.Errorf("Should have no issues, found %d", len(report.Issues))
	}
}

func TestCompileScenes_MixedIssues(t *testing.T) {
	scenes := map[string]*Scene{
		// Valid scene
		"valid": {ID: "valid", Chapter: 1, SceneNumber: 1},
		// Orphan
		"orphan": {ID: "orphan", Name: "Orphan", Chapter: 0, SceneNumber: 0},
		// Gap (missing scene 2)
		"scene3": {ID: "scene3", Chapter: 1, SceneNumber: 3},
		// Duplicates
		"dup1": {ID: "dup1", Chapter: 2, SceneNumber: 1},
		"dup2": {ID: "dup2", Chapter: 2, SceneNumber: 1},
	}

	report := CompileScenes(scenes)

	if len(report.Orphans) != 1 {
		t.Errorf("Orphans = %d, expected 1", len(report.Orphans))
	}

	if len(report.Gaps) != 1 {
		t.Errorf("Gaps = %d, expected 1", len(report.Gaps))
	}

	if len(report.Duplicates) != 1 {
		t.Errorf("Duplicates = %d, expected 1", len(report.Duplicates))
	}

	// Should have: 1 orphan warning, 1 gap warning, 1 duplicate error
	if len(report.Issues) != 3 {
		t.Errorf("Total issues = %d, expected 3", len(report.Issues))
	}

	// Issues should be sorted by severity (errors first)
	if len(report.Issues) > 0 && report.Issues[0].Severity != "error" {
		t.Error("First issue should be an error (duplicates)")
	}
}

func TestCompilationReport_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		issues   []CompilationIssue
		expected bool
	}{
		{
			name:     "no issues",
			issues:   []CompilationIssue{},
			expected: false,
		},
		{
			name: "only warnings",
			issues: []CompilationIssue{
				{Severity: "warning"},
				{Severity: "warning"},
			},
			expected: false,
		},
		{
			name: "has error",
			issues: []CompilationIssue{
				{Severity: "warning"},
				{Severity: "error"},
			},
			expected: true,
		},
		{
			name: "only errors",
			issues: []CompilationIssue{
				{Severity: "error"},
				{Severity: "error"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &CompilationReport{Issues: tt.issues}
			result := report.HasErrors()
			if result != tt.expected {
				t.Errorf("HasErrors() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCompilationReport_HasWarnings(t *testing.T) {
	tests := []struct {
		name     string
		issues   []CompilationIssue
		expected bool
	}{
		{
			name:     "no issues",
			issues:   []CompilationIssue{},
			expected: false,
		},
		{
			name: "only errors",
			issues: []CompilationIssue{
				{Severity: "error"},
				{Severity: "error"},
			},
			expected: false,
		},
		{
			name: "has warning",
			issues: []CompilationIssue{
				{Severity: "error"},
				{Severity: "warning"},
			},
			expected: true,
		},
		{
			name: "only warnings",
			issues: []CompilationIssue{
				{Severity: "warning"},
				{Severity: "warning"},
			},
			expected: true,
		},
		{
			name: "info only",
			issues: []CompilationIssue{
				{Severity: "info"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &CompilationReport{Issues: tt.issues}
			result := report.HasWarnings()
			if result != tt.expected {
				t.Errorf("HasWarnings() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCompilationReport_Summary(t *testing.T) {
	tests := []struct {
		name     string
		issues   []CompilationIssue
		expected string
	}{
		{
			name:     "no issues",
			issues:   []CompilationIssue{},
			expected: "✓ All scenes are properly organized!",
		},
		{
			name: "one error",
			issues: []CompilationIssue{
				{Severity: "error"},
			},
			expected: "Found 1 error(s) and 0 warning(s)",
		},
		{
			name: "one warning",
			issues: []CompilationIssue{
				{Severity: "warning"},
			},
			expected: "Found 0 error(s) and 1 warning(s)",
		},
		{
			name: "multiple errors and warnings",
			issues: []CompilationIssue{
				{Severity: "error"},
				{Severity: "error"},
				{Severity: "warning"},
				{Severity: "warning"},
				{Severity: "warning"},
			},
			expected: "Found 2 error(s) and 3 warning(s)",
		},
		{
			name: "with info",
			issues: []CompilationIssue{
				{Severity: "error"},
				{Severity: "warning"},
				{Severity: "info"},
			},
			expected: "Found 1 error(s) and 1 warning(s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &CompilationReport{Issues: tt.issues}
			result := report.Summary()
			if result != tt.expected {
				t.Errorf("Summary() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestCompileScenes_IssueSorting(t *testing.T) {
	scenes := map[string]*Scene{
		"orphan": {ID: "orphan", Name: "Orphan", Chapter: 0, SceneNumber: 0},
		"gap1":   {ID: "gap1", Chapter: 1, SceneNumber: 1},
		"gap2":   {ID: "gap2", Chapter: 1, SceneNumber: 3},
		"dup1":   {ID: "dup1", Chapter: 2, SceneNumber: 1},
		"dup2":   {ID: "dup2", Chapter: 2, SceneNumber: 1},
	}

	report := CompileScenes(scenes)

	// First issue should be error (duplicate)
	if len(report.Issues) > 0 && report.Issues[0].Severity != "error" {
		t.Error("First issue should be an error")
	}

	// Count errors at the beginning
	errorCount := 0
	for i, issue := range report.Issues {
		if issue.Severity == "error" {
			errorCount++
		} else {
			// Once we hit non-error, all remaining should be non-errors
			for j := i + 1; j < len(report.Issues); j++ {
				if report.Issues[j].Severity == "error" {
					t.Error("Errors should be sorted before warnings")
				}
			}
			break
		}
	}
}
