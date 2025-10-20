package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatExport_BPMAndChordProKeyDetection(t *testing.T) {
	ef := NewExportFormatter()

	content := `# Demo Song
Tempo: 128

[C] This is a chord line
G Am F C
`

	opts := DefaultExportOptions()
	opts.Type = ExportTypeChordPro
	opts.Title = "Demo Song"
	opts.BPM = 0 // let formatter extract BPM

	ne, err := ef.FormatExport(content, opts)
	if err != nil {
		t.Fatalf("FormatExport failed: %v", err)
	}

	if !strings.Contains(ne.Lyrics, "{tempo:128}") && !strings.Contains(ne.Lyrics, "{tempo: 128}") {
		t.Errorf("expected tempo directive in chordpro output, got: %q", ne.Lyrics)
	}

	// detectKey should pick up a root (C or G depending on content). Ensure a key directive exists.
	if !strings.Contains(ne.Lyrics, "{key:") {
		t.Errorf("expected key directive in chordpro output, got: %q", ne.Lyrics)
	}
}

func TestExportFormatter_SaveToFile_DifferentExtensions(t *testing.T) {
	ef := NewExportFormatter()
	tmp := t.TempDir()

	ne := &NoiseExport{
		Type:    "noise_sh_test",
		Version: "test",
		Lyrics:  "line1\nline2",
		Metadata: ExportMetadata{
			Title: "T",
			BPM:   90,
		},
	}

	// JSON
	jsonPath := filepath.Join(tmp, "out.json")
	if err := ef.SaveToFile(ne, jsonPath); err != nil {
		t.Fatalf("SaveToFile json failed: %v", err)
	}
	b, _ := os.ReadFile(jsonPath)
	var parsed NoiseExport
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Fatalf("invalid json written: %v", err)
	}
	if parsed.Type != ne.Type {
		t.Errorf("json content mismatch: got type %q want %q", parsed.Type, ne.Type)
	}

	// Markdown (.md) should write Lyrics directly
	mdPath := filepath.Join(tmp, "out.md")
	if err := ef.SaveToFile(ne, mdPath); err != nil {
		t.Fatalf("SaveToFile md failed: %v", err)
	}
	b, _ = os.ReadFile(mdPath)
	if string(b) != ne.Lyrics {
		t.Errorf(".md write mismatch: got %q want %q", string(b), ne.Lyrics)
	}

	// Plain text (.txt)
	txtPath := filepath.Join(tmp, "out.txt")
	if err := ef.SaveToFile(ne, txtPath); err != nil {
		t.Fatalf("SaveToFile txt failed: %v", err)
	}
	b, _ = os.ReadFile(txtPath)
	if string(b) != ne.Lyrics {
		t.Errorf(".txt write mismatch: got %q want %q", string(b), ne.Lyrics)
	}

	// ChordPro (.cho)
	choPath := filepath.Join(tmp, "out.cho")
	if err := ef.SaveToFile(ne, choPath); err != nil {
		t.Fatalf("SaveToFile cho failed: %v", err)
	}
	b, _ = os.ReadFile(choPath)
	if string(b) != ne.Lyrics {
		t.Errorf(".cho write mismatch: got %q want %q", string(b), ne.Lyrics)
	}
}

func TestExportService_ExportAndHelpers(t *testing.T) {
	outDir := t.TempDir()
	es := NewExportService(outDir)

	content := `# Simple
BPM: 95

[C]hello`

	// Export as plain text to a specific filename
	opts := DefaultExportOptions()
	opts.Type = ExportTypePlainText
	opts.OutputPath = "my_export.txt"

	path, err := es.Export(content, opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Path should be placed inside outputDir
	if !strings.HasPrefix(path, outDir) {
		t.Fatalf("expected output path to be inside %s, got %s", outDir, path)
	}

	// File exists
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read exported file: %v", err)
	}

	// Plain text formatting should remove markdown-like markers such as "#"
	if strings.Contains(string(data), "#") {
		t.Errorf("plain text export contains markdown header: %q", string(data))
	}

	// ListExports should find JSON files; create two json files
	os.WriteFile(filepath.Join(outDir, "a.json"), []byte("{}"), 0644)
	os.WriteFile(filepath.Join(outDir, "b.json"), []byte("{}"), 0644)

	exports, err := es.ListExports()
	if err != nil {
		t.Fatalf("ListExports failed: %v", err)
	}
	// Expect at least the two json files (order not guaranteed)
	foundA, foundB := false, false
	for _, e := range exports {
		if e == "a.json" {
			foundA = true
		}
		if e == "b.json" {
			foundB = true
		}
	}
	if !foundA || !foundB {
		t.Errorf("ListExports missing created files: got %v", exports)
	}

	// DeleteExport should remove a file
	if err := es.DeleteExport("a.json"); err != nil {
		t.Fatalf("DeleteExport failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "a.json")); !os.IsNotExist(err) {
		t.Errorf("expected a.json to be deleted, but it exists or stat failed: %v", err)
	}

	// GetExportPath with relative name
	got := es.GetExportPath("b.json")
	if !strings.HasSuffix(got, "b.json") {
		t.Errorf("GetExportPath unexpected result: %s", got)
	}
}

func TestExportService_ErrorOnWriteToDirectory(t *testing.T) {
	outDir := t.TempDir()
	es := NewExportService(outDir)

	// Force a write error by setting OutputPath to "." which will cause the service
	// to attempt to write to the output directory path itself (a directory), leading to a write error.
	opts := DefaultExportOptions()
	opts.OutputPath = "." // will be joined into outDir/. which resolves to a directory

	_, err := es.Export("content", opts)
	if err == nil {
		t.Fatalf("expected an error when writing to a directory path, got nil")
	}
	// Ensure error message indicates a write/save problem
	if !strings.Contains(err.Error(), "save") && !strings.Contains(err.Error(), "write") {
		t.Logf("received error (ok) but unexpected message: %v", err)
	}
}
