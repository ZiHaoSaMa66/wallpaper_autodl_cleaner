//go:build windows

package steam

import (
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
	"wp-cleaner/model"
)

func FindSteamPath() (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("cannot open Steam registry key: %w", err)
	}
	defer k.Close()

	steamPath, _, err := k.GetStringValue("SteamPath")
	if err != nil {
		return "", fmt.Errorf("cannot read SteamPath from registry: %w", err)
	}
	// Steam registry key uses forward slashes
	steamPath = filepath.FromSlash(steamPath)
	return steamPath, nil
}

func GetWorkshopPath(steamPath string) string {
	return filepath.Join(steamPath, "steamapps", "workshop", "content", model.WallpaperEngineAppID)
}

func GetLoginUsersPath(steamPath string) string {
	return filepath.Join(steamPath, "config", "loginusers.vdf")
}

// FindWEPath returns the Wallpaper Engine install path.
func FindWEPath(steamPath string) string {
	return filepath.Join(steamPath, "steamapps", "common", "wallpaper_engine")
}

func CheckProcessesRunning(names ...string) ([]string, bool) {
	if len(names) == 0 {
		names = []string{"steam.exe"}
	}
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: cannot check running processes: %v\n", err)
		return nil, false
	}
	defer windows.CloseHandle(snapshot)

	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	err = windows.Process32First(snapshot, &pe)
	if err != nil {
		return nil, false
	}

	var running []string
	for err == nil {
		exeName := windows.UTF16ToString(pe.ExeFile[:])
		for _, n := range names {
			if strings.EqualFold(exeName, n) {
				running = append(running, exeName)
				break
			}
		}
		err = windows.Process32Next(snapshot, &pe)
	}
	return running, len(running) > 0
}

func SteamIsRunning() bool {
	_, running := CheckProcessesRunning("steam.exe")
	return running
}
