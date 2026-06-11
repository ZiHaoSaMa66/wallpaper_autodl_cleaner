package steam

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jslay88/vdf"
	"wp-cleaner/model"
)

// ReadLibraryFoldersVDF reads Steam's libraryfolders.vdf and returns all library paths.
func ReadLibraryFoldersVDF(steamPath string) ([]string, error) {
	vdfPath := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
	data, err := os.ReadFile(vdfPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", vdfPath, err)
	}

	doc, err := vdf.ParseString(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing libraryfolders.vdf: %w", err)
	}

	root := doc.Get("libraryfolders")
	if root == nil {
		return nil, nil
	}

	var libraries []string
	for _, child := range root.Children {
		if !child.IsObject {
			continue
		}
		path := child.GetString("path")
		if path != "" {
			libraries = append(libraries, path)
		}
	}
	return libraries, nil
}

// FindWELibraryPath finds the Steam library folder where Wallpaper Engine is installed.
// It first checks libraryfolders.vdf for the app ID mapping, then falls back to
// checking folder existence in all libraries including the main Steam path.
func FindWELibraryPath(steamPath string) (string, error) {
	// First try to find via libraryfolders.vdf app mapping (most reliable)
	vdfPath := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
	if data, err := os.ReadFile(vdfPath); err == nil {
		if doc, err := vdf.ParseString(string(data)); err == nil {
			root := doc.Get("libraryfolders")
			if root != nil {
				for _, child := range root.Children {
					if !child.IsObject {
						continue
					}
					appsNode := child.Get("apps")
					if appsNode != nil {
						for _, appChild := range appsNode.Children {
							if appChild.Key == model.WallpaperEngineAppID {
								libPath := child.GetString("path")
								if libPath != "" {
									return libPath, nil
								}
							}
						}
					}
				}
			}
		}
	}

	// Fallback: check main Steam path for workshop or common folder
	wsPath := filepath.Join(steamPath, "steamapps", "workshop", "content", model.WallpaperEngineAppID)
	if _, err := os.Stat(wsPath); err == nil {
		return steamPath, nil
	}
	commonPath := filepath.Join(steamPath, "steamapps", "common", "wallpaper_engine")
	if _, err := os.Stat(commonPath); err == nil {
		return steamPath, nil
	}

	// Check all secondary libraries
	libraries, err := ReadLibraryFoldersVDF(steamPath)
	if err != nil {
		return "", fmt.Errorf("wallpaper engine not found in main Steam path and cannot read libraryfolders.vdf: %w", err)
	}

	for _, libPath := range libraries {
		if libPath == steamPath {
			continue
		}
		wsPath := filepath.Join(libPath, "steamapps", "workshop", "content", model.WallpaperEngineAppID)
		if _, err := os.Stat(wsPath); err == nil {
			return libPath, nil
		}
		commonPath := filepath.Join(libPath, "steamapps", "common", "wallpaper_engine")
		if _, err := os.Stat(commonPath); err == nil {
			return libPath, nil
		}
	}

	return "", fmt.Errorf("wallpaper engine not found in any Steam library folder")
}

// GetWEDownloadsPath returns the Wallpaper Engine workshop downloads path.
func GetWEDownloadsPath(libraryPath string) string {
	return filepath.Join(libraryPath, "steamapps", "workshop", "downloads")
}

// GetWEACFPath returns the Wallpaper Engine workshop ACF file path.
func GetWEACFPath(libraryPath string) string {
	return filepath.Join(libraryPath, "steamapps", "workshop", "appworkshop_"+model.WallpaperEngineAppID+".acf")
}
