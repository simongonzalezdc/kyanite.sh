package scene

import "sort"

// SortScenes converts a map of scenes to a sorted slice by chapter and scene number
func SortScenes(scenes map[string]*Scene) []*Scene {
	sorted := make([]*Scene, 0, len(scenes))
	for _, sc := range scenes {
		sorted = append(sorted, sc)
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Chapter != sorted[j].Chapter {
			return sorted[i].Chapter < sorted[j].Chapter
		}
		return sorted[i].SceneNumber < sorted[j].SceneNumber
	})

	return sorted
}
