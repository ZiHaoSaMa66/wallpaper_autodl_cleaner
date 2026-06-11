package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func ScanWorkshop(workshopPath string) ([]uint64, error) {
	entries, err := os.ReadDir(workshopPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read workshop dir %s: %w", workshopPath, err)
	}
	var ids []uint64
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		id, err := strconv.ParseUint(e.Name(), 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// GetWorkshopDirSize returns the total size of all files in the workshop directory.
func GetWorkshopDirSize(workshopPath string) (int64, error) {
	var total int64
	err := filepath.WalkDir(workshopPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				total += info.Size()
			}
		}
		return nil
	})
	return total, err
}

func HumanSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}


