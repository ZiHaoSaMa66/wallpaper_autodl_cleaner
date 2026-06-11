package cleaner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanTrash_Empty(t *testing.T) {
	dir := t.TempDir()

	paths, size, err := ScanTrash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 0 {
		t.Errorf("expected 0 paths, got %d", len(paths))
	}
	if size != 0 {
		t.Errorf("expected 0 size, got %d", size)
	}
}

func TestScanTrash_FindsTrashFolders(t *testing.T) {
	dir := t.TempDir()

	trashDir := filepath.Join(dir, ".trash-12345")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		t.Fatal(err)
	}
	file1 := filepath.Join(trashDir, "file1.bin")
	if err := os.WriteFile(file1, make([]byte, 100), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(trashDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	file2 := filepath.Join(subDir, "file2.bin")
	if err := os.WriteFile(file2, make([]byte, 50), 0644); err != nil {
		t.Fatal(err)
	}

	paths, size, err := ScanTrash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d: %v", len(paths), paths)
	}
	if !strings.HasSuffix(paths[0], ".trash-12345") {
		t.Errorf("expected path to end with .trash-12345, got %s", paths[0])
	}
	if size != 150 {
		t.Errorf("expected size 150, got %d", size)
	}
}

func TestScanTrash_IgnoresNonTrash(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "12345"), 0755)
	os.MkdirAll(filepath.Join(dir, "67890"), 0755)
	os.WriteFile(filepath.Join(dir, "12345", "a.bin"), make([]byte, 10), 0644)

	paths, size, err := ScanTrash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 0 {
		t.Errorf("expected 0 paths for non-trash dirs, got %d", len(paths))
	}
	if size != 0 {
		t.Errorf("expected 0 size, got %d", size)
	}
}
