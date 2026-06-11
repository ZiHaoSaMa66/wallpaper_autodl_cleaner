package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"wp-cleaner/model"
	"wp-cleaner/scanner"
	"wp-cleaner/steam"
)

type CleanupPlan struct {
	ToRemove     []model.WallpaperInfo
	ToKeep       []model.WallpaperInfo
	TotalSize    int64
	WorkshopPath string
}

func BuildCleanupPlan(
	workshopPath string,
	localIDs []uint64,
	subscribedIDs []uint64,
	wallpaperInfos []model.WallpaperInfo,
) *CleanupPlan {
	subscribed := steam.BuildSubscribedSet(subscribedIDs)

	infoMap := make(map[uint64]model.WallpaperInfo)
	for _, info := range wallpaperInfos {
		info.LocalPath = filepath.Join(workshopPath, fmt.Sprintf("%d", info.PublishedFileID))
		infoMap[info.PublishedFileID] = info
	}

	plan := &CleanupPlan{WorkshopPath: workshopPath}

	for _, id := range localIDs {
		info, known := infoMap[id]
		if !known {
			info = model.WallpaperInfo{
				PublishedFileID: id,
				Title:           fmt.Sprintf("Unknown (ID: %d)", id),
				LocalPath:       filepath.Join(workshopPath, fmt.Sprintf("%d", id)),
			}
		}
		if subscribed[id] {
			plan.ToKeep = append(plan.ToKeep, info)
		} else {
			plan.ToRemove = append(plan.ToRemove, info)
			plan.TotalSize += int64(info.FileSize)
		}
	}
	return plan
}

func (p *CleanupPlan) DryRun() {
	fmt.Printf("Workshop Path: %s\n", p.WorkshopPath)
	fmt.Printf("Total wallpapers on disk: %d\n", len(p.ToRemove)+len(p.ToKeep))
	fmt.Printf("Subscribed (will keep):  %d\n", len(p.ToKeep))
	fmt.Printf("Unsubscribed (to remove): %d\n", len(p.ToRemove))
	if len(p.ToRemove) > 0 {
		fmt.Printf("Estimated space to free: %s\n", scanner.HumanSize(p.TotalSize))
		fmt.Println("\nWallpapers to remove:")
		for _, w := range p.ToRemove {
			fmt.Printf("  [%d] %s\n", w.PublishedFileID, w.Title)
		}
	}
}

func (p *CleanupPlan) Execute() error {
	if len(p.ToRemove) == 0 {
		fmt.Println("Nothing to clean.")
		return nil
	}

	fmt.Printf("Quarantining %d wallpapers (moving to .trash-* folders)...\n", len(p.ToRemove))
	var errs int
	for _, w := range p.ToRemove {
		trashPath := filepath.Join(p.WorkshopPath, fmt.Sprintf(".trash-%d", w.PublishedFileID))
		err := os.Rename(w.LocalPath, trashPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: failed to move %d: %v\n", w.PublishedFileID, err)
			errs++
		} else {
			fmt.Printf("  Moved: %s (ID: %d)\n", w.Title, w.PublishedFileID)
		}
	}
	fmt.Println("Done. Review quarantined folders in the workshop directory.")
	fmt.Println("To restore, rename folders from .trash-* back to original IDs.")
	if errs > 0 {
		return fmt.Errorf("%d items failed to quarantine", errs)
	}
	return nil
}

func ScanTrash(workshopPath string) (paths []string, totalSize int64, err error) {
	entries, err := os.ReadDir(workshopPath)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot read workshop dir: %w", err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if !strings.HasPrefix(e.Name(), ".trash-") {
			continue
		}
		fullPath := filepath.Join(workshopPath, e.Name())
		var size int64
		err = filepath.WalkDir(fullPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				info, err := d.Info()
				if err == nil {
					size += info.Size()
				}
			}
			return nil
		})
		if err != nil {
			return nil, 0, fmt.Errorf("cannot calculate size for %s: %w", e.Name(), err)
		}
		paths = append(paths, fullPath)
		totalSize += size
	}
	return paths, totalSize, nil
}

func DeleteTrash(workshopPath string, dryRun bool, force bool) error {
	paths, totalSize, err := ScanTrash(workshopPath)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		fmt.Println("No trash folders found.")
		return nil
	}

	fmt.Printf("Found %d trash folder(s), total size: %s\n", len(paths), scanner.HumanSize(totalSize))
	for _, p := range paths {
		fmt.Printf("  %s\n", p)
	}

	if dryRun {
		fmt.Println("Dry-run mode: nothing was deleted.")
		return nil
	}

	if !force {
		fmt.Print("Permanently delete these folders? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	var errs int
	for _, p := range paths {
		if err := os.RemoveAll(p); err != nil {
			fmt.Fprintf(os.Stderr, "WARN: failed to delete %s: %v\n", p, err)
			errs++
		} else {
			fmt.Printf("  Deleted: %s\n", p)
		}
	}
	if errs > 0 {
		return fmt.Errorf("%d items failed to delete", errs)
	}
	return nil
}
