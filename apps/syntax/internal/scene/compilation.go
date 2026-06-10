package scene

import (
	"fmt"
	"sort"
)

// CompilationIssue represents a problem found during scene compilation
type CompilationIssue struct {
	Type        string
	Severity    string // "error", "warning", "info"
	Message     string
	SceneID     string
	Chapter     int
	SceneNumber int
}

// CompilationReport contains the results of scene analysis
type CompilationReport struct {
	TotalScenes  int
	Issues       []CompilationIssue
	Orphans      []*Scene
	Gaps         []GapInfo
	Duplicates   []DuplicateInfo
	ChapterCount int
}

// GapInfo represents a missing scene in the sequence
type GapInfo struct {
	Chapter     int
	SceneNumber int
	Message     string
}

// DuplicateInfo represents duplicate scene numbers
type DuplicateInfo struct {
	Chapter     int
	SceneNumber int
	SceneIDs    []string
	Message     string
}

// CompileScenes analyzes scenes for issues and returns a report
func CompileScenes(scenes map[string]*Scene) *CompilationReport {
	report := &CompilationReport{
		TotalScenes: len(scenes),
		Issues:      []CompilationIssue{},
		Orphans:     []*Scene{},
		Gaps:        []GapInfo{},
		Duplicates:  []DuplicateInfo{},
	}

	if len(scenes) == 0 {
		report.Issues = append(report.Issues, CompilationIssue{
			Type:     "empty",
			Severity: "info",
			Message:  "No scenes to compile",
		})
		return report
	}

	// Convert to sorted list
	sceneList := make([]*Scene, 0, len(scenes))
	for _, sc := range scenes {
		sceneList = append(sceneList, sc)
	}
	sort.Slice(sceneList, func(i, j int) bool {
		if sceneList[i].Chapter != sceneList[j].Chapter {
			return sceneList[i].Chapter < sceneList[j].Chapter
		}
		return sceneList[i].SceneNumber < sceneList[j].SceneNumber
	})

	// Detect orphans (scenes without chapter or scene number)
	for _, sc := range sceneList {
		if sc.Chapter == 0 || sc.SceneNumber == 0 {
			report.Orphans = append(report.Orphans, sc)
			report.Issues = append(report.Issues, CompilationIssue{
				Type:        "orphan",
				Severity:    "warning",
				Message:     fmt.Sprintf("Scene '%s' has no chapter/scene number", sc.Name),
				SceneID:     sc.ID,
				Chapter:     sc.Chapter,
				SceneNumber: sc.SceneNumber,
			})
		}
	}

	// Group scenes by chapter
	chapterMap := make(map[int][]*Scene)
	for _, sc := range sceneList {
		if sc.Chapter > 0 {
			chapterMap[sc.Chapter] = append(chapterMap[sc.Chapter], sc)
		}
	}

	report.ChapterCount = len(chapterMap)

	// Check each chapter for gaps and duplicates
	for chapter, chapterScenes := range chapterMap {
		// Check for duplicates
		sceneNumMap := make(map[int][]string)
		for _, sc := range chapterScenes {
			sceneNumMap[sc.SceneNumber] = append(sceneNumMap[sc.SceneNumber], sc.ID)
		}

		for sceneNum, ids := range sceneNumMap {
			if len(ids) > 1 {
				dup := DuplicateInfo{
					Chapter:     chapter,
					SceneNumber: sceneNum,
					SceneIDs:    ids,
					Message:     fmt.Sprintf("Chapter %d has duplicate scene %d", chapter, sceneNum),
				}
				report.Duplicates = append(report.Duplicates, dup)
				report.Issues = append(report.Issues, CompilationIssue{
					Type:        "duplicate",
					Severity:    "error",
					Message:     dup.Message,
					Chapter:     chapter,
					SceneNumber: sceneNum,
				})
			}
		}

		// Check for gaps
		sort.Slice(chapterScenes, func(i, j int) bool {
			return chapterScenes[i].SceneNumber < chapterScenes[j].SceneNumber
		})

		for i := 0; i < len(chapterScenes)-1; i++ {
			current := chapterScenes[i].SceneNumber
			next := chapterScenes[i+1].SceneNumber

			if next-current > 1 {
				// There's a gap
				for missing := current + 1; missing < next; missing++ {
					gap := GapInfo{
						Chapter:     chapter,
						SceneNumber: missing,
						Message:     fmt.Sprintf("Chapter %d is missing scene %d", chapter, missing),
					}
					report.Gaps = append(report.Gaps, gap)
					report.Issues = append(report.Issues, CompilationIssue{
						Type:        "gap",
						Severity:    "warning",
						Message:     gap.Message,
						Chapter:     chapter,
						SceneNumber: missing,
					})
				}
			}
		}

		// Check if chapter starts at scene 1
		if len(chapterScenes) > 0 && chapterScenes[0].SceneNumber > 1 {
			for missing := 1; missing < chapterScenes[0].SceneNumber; missing++ {
				gap := GapInfo{
					Chapter:     chapter,
					SceneNumber: missing,
					Message:     fmt.Sprintf("Chapter %d is missing scene %d at start", chapter, missing),
				}
				report.Gaps = append(report.Gaps, gap)
				report.Issues = append(report.Issues, CompilationIssue{
					Type:        "gap",
					Severity:    "warning",
					Message:     gap.Message,
					Chapter:     chapter,
					SceneNumber: missing,
				})
			}
		}
	}

	// Sort issues by severity (errors first)
	sort.Slice(report.Issues, func(i, j int) bool {
		severityOrder := map[string]int{"error": 0, "warning": 1, "info": 2}
		return severityOrder[report.Issues[i].Severity] < severityOrder[report.Issues[j].Severity]
	})

	return report
}

// HasErrors returns true if the report contains any errors
func (r *CompilationReport) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

// HasWarnings returns true if the report contains any warnings
func (r *CompilationReport) HasWarnings() bool {
	for _, issue := range r.Issues {
		if issue.Severity == "warning" {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary
func (r *CompilationReport) Summary() string {
	if len(r.Issues) == 0 {
		return "✓ All scenes are properly organized!"
	}

	errorCount := 0
	warningCount := 0
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			errorCount++
		} else if issue.Severity == "warning" {
			warningCount++
		}
	}

	return fmt.Sprintf("Found %d error(s) and %d warning(s)", errorCount, warningCount)
}
