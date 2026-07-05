package appnames

import (
	"slices"
	"testing"
)

func TestAllContainsFourApps(t *testing.T) {
	if len(All) != 4 {
		t.Fatalf("expected 4 apps, got %d", len(All))
	}
}

func TestOthersExcludesSelf(t *testing.T) {
	for _, app := range All {
		others := Others(app)
		if len(others) != 3 {
			t.Errorf("Others(%q): expected 3, got %d", app, len(others))
		}
		if slices.Contains(others, app) {
			t.Errorf("Others(%q): should not contain self", app)
		}
	}
}

func TestOthersIsStable(t *testing.T) {
	a := Others(Focus)
	b := Others(Focus)
	if len(a) != len(b) {
		t.Fatal("Others should return consistent length")
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("Others not stable: %v vs %v", a, b)
		}
	}
}
