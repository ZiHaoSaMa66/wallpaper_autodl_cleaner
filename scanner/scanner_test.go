package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHumanSize_Bytes(t *testing.T) {
	if s := HumanSize(0); s != "0 B" {
		t.Errorf("expected 0 B, got %s", s)
	}
	if s := HumanSize(1); s != "1 B" {
		t.Errorf("expected 1 B, got %s", s)
	}
	if s := HumanSize(1023); s != "1023 B" {
		t.Errorf("expected 1023 B, got %s", s)
	}
}

func TestHumanSize_KB(t *testing.T) {
	if s := HumanSize(1024); s != "1.0 KB" {
		t.Errorf("expected 1.0 KB, got %s", s)
	}
	if s := HumanSize(2048); s != "2.0 KB" {
		t.Errorf("expected 2.0 KB, got %s", s)
	}
	if s := HumanSize(1536); s != "1.5 KB" {
		t.Errorf("expected 1.5 KB, got %s", s)
	}
}

func TestHumanSize_MB(t *testing.T) {
	if s := HumanSize(1024 * 1024); s != "1.0 MB" {
		t.Errorf("expected 1.0 MB, got %s", s)
	}
	if s := HumanSize(5*1024*1024 + 512*1024); s != "5.5 MB" {
		t.Errorf("expected 5.5 MB, got %s", s)
	}
}

func TestHumanSize_GB(t *testing.T) {
	if s := HumanSize(1024 * 1024 * 1024); s != "1.0 GB" {
		t.Errorf("expected 1.0 GB, got %s", s)
	}
}

func TestHumanSize_TB(t *testing.T) {
	if s := HumanSize(1024 * 1024 * 1024 * 1024); s != "1.0 TB" {
		t.Errorf("expected 1.0 TB, got %s", s)
	}
}

func TestScanWorkshop_ValidIDs(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "12345"), 0755)
	os.MkdirAll(filepath.Join(dir, "67890"), 0755)

	ids, err := ScanWorkshop(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %d: %v", len(ids), ids)
	}
	if ids[0] != 12345 && ids[0] != 67890 {
		t.Errorf("unexpected id: %d", ids[0])
	}
	if ids[1] != 12345 && ids[1] != 67890 {
		t.Errorf("unexpected id: %d", ids[1])
	}
}

func TestScanWorkshop_IgnoresNonNumeric(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "abc"), 0755)
	os.MkdirAll(filepath.Join(dir, "123"), 0755)

	ids, err := ScanWorkshop(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 {
		t.Fatalf("expected 1 id, got %d: %v", len(ids), ids)
	}
	if ids[0] != 123 {
		t.Errorf("expected 123, got %d", ids[0])
	}
}

func TestScanWorkshop_IgnoresFiles(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "123"), 0755)
	os.WriteFile(filepath.Join(dir, "notadir"), []byte("test"), 0644)

	ids, err := ScanWorkshop(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 {
		t.Fatalf("expected 1 id, got %d: %v", len(ids), ids)
	}
	if ids[0] != 123 {
		t.Errorf("expected 123, got %d", ids[0])
	}
}

func TestScanWorkshop_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	ids, err := ScanWorkshop(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 0 {
		t.Errorf("expected 0 ids, got %d", len(ids))
	}
}
