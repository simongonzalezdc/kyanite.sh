package export

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRegistryRegisterAndLookup(t *testing.T) {
	r := NewRegistry()
	f := NewFormatter("css", func(data any) ([]byte, error) {
		return []byte("formatted"), nil
	})
	r.Register("palette", "css", f)

	found := r.Lookup("palette", "css")
	if found == nil {
		t.Fatal("expected to find registered formatter")
	}

	data, err := found.Format(nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "formatted" {
		t.Errorf("got %q, want %q", data, "formatted")
	}
	if found.Extension() != "css" {
		t.Errorf("Extension() = %q, want %q", found.Extension(), "css")
	}
}

func TestRegistryLookupMissing(t *testing.T) {
	r := NewRegistry()
	if r.Lookup("missing", "missing") != nil {
		t.Error("expected nil for unregistered formatter")
	}
}

func TestRegistryFormats(t *testing.T) {
	r := NewRegistry()
	r.Register("palette", "css", NewFormatter("css", func(any) ([]byte, error) { return nil, nil }))
	r.Register("palette", "json", NewFormatter("json", func(any) ([]byte, error) { return nil, nil }))

	formats := r.Formats("palette")
	if len(formats) != 2 {
		t.Fatalf("expected 2 formats, got %d", len(formats))
	}
}

func TestExportUnknownFormat(t *testing.T) {
	r := NewRegistry()
	_, err := Export(r, "palette", "xml", nil)
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestWriteFileCreatesDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "test.txt")
	if err := WriteFile(path, []byte("hello")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Errorf("got %q, want %q", data, "hello")
	}
}
