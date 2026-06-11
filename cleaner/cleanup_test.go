package cleaner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"wp-cleaner/model"
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

func TestBuildCleanupPlan_Basic(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "111"), 0755)
	os.MkdirAll(filepath.Join(dir, "222"), 0755)

	localIDs := []uint64{111, 222}
	subscribedIDs := []uint64{111}
	infos := []model.WallpaperInfo{
		{PublishedFileID: 111, Title: "Subbed WP", FileSize: 100},
		{PublishedFileID: 222, Title: "Unsubbed WP", FileSize: 200},
	}

	plan := BuildCleanupPlan(dir, localIDs, subscribedIDs, infos)

	if len(plan.ToKeep) != 1 || plan.ToKeep[0].PublishedFileID != 111 {
		t.Errorf("expected 1 kept item (111), got %d items", len(plan.ToKeep))
	}
	if len(plan.ToRemove) != 1 || plan.ToRemove[0].PublishedFileID != 222 {
		t.Errorf("expected 1 removed item (222), got %d items", len(plan.ToRemove))
	}
	if plan.TotalSize != 200 {
		t.Errorf("expected total size 200, got %d", plan.TotalSize)
	}
}

func TestBuildCleanupPlan_UnknownID(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "999"), 0755)

	localIDs := []uint64{999}
	plan := BuildCleanupPlan(dir, localIDs, nil, nil)

	if len(plan.ToRemove) != 1 {
		t.Fatalf("expected 1 item to remove, got %d", len(plan.ToRemove))
	}
	if plan.ToRemove[0].Title != "Unknown (ID: 999)" {
		t.Errorf("expected 'Unknown (ID: 999)', got '%s'", plan.ToRemove[0].Title)
	}
}

func TestBuildCleanupPlan_NoSubscriptions(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "111"), 0755)
	os.MkdirAll(filepath.Join(dir, "222"), 0755)

	localIDs := []uint64{111, 222}
	plan := BuildCleanupPlan(dir, localIDs, nil, nil)

	if len(plan.ToKeep) != 0 {
		t.Errorf("expected 0 kept items, got %d", len(plan.ToKeep))
	}
	if len(plan.ToRemove) != 2 {
		t.Errorf("expected 2 removed items, got %d", len(plan.ToRemove))
	}
}

func TestDeleteTrash_DryRun(t *testing.T) {
	dir := t.TempDir()
	trashDir := filepath.Join(dir, ".trash-111")
	os.MkdirAll(trashDir, 0755)
	os.WriteFile(filepath.Join(trashDir, "f.bin"), make([]byte, 50), 0644)

	err := DeleteTrash(dir, true, true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(trashDir); os.IsNotExist(err) {
		t.Error("trash dir was deleted despite dry-run mode")
	}
}

func TestDeleteTrash_Force(t *testing.T) {
	dir := t.TempDir()
	trashDir := filepath.Join(dir, ".trash-111")
	os.MkdirAll(trashDir, 0755)
	os.WriteFile(filepath.Join(trashDir, "f.bin"), make([]byte, 50), 0644)

	err := DeleteTrash(dir, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(trashDir); !os.IsNotExist(err) {
		t.Error("trash dir should have been deleted")
	}
}

func TestDeleteTrash_NoTrash(t *testing.T) {
	dir := t.TempDir()
	err := DeleteTrash(dir, false, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecute_MovesToTrash(t *testing.T) {
	dir := t.TempDir()

	wpDir := filepath.Join(dir, "111")
	os.MkdirAll(wpDir, 0755)
	os.WriteFile(filepath.Join(wpDir, "a.bin"), make([]byte, 10), 0644)

	plan := &CleanupPlan{
		WorkshopPath: dir,
		ToRemove: []model.WallpaperInfo{
			{PublishedFileID: 111, Title: "Test WP", LocalPath: wpDir},
		},
	}

	err := plan.Execute()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(wpDir); !os.IsNotExist(err) {
		t.Error("original wallpaper dir should no longer exist")
	}

	trashDir := filepath.Join(dir, ".trash-111")
	if _, err := os.Stat(trashDir); os.IsNotExist(err) {
		t.Error("trash folder should exist after Execute")
	}
}

func TestExecute_NothingToRemove(t *testing.T) {
	dir := t.TempDir()
	plan := &CleanupPlan{
		WorkshopPath: dir,
		ToRemove:     nil,
	}

	err := plan.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
