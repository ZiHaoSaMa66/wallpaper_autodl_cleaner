package steam

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
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
	return filepath.Join(steamPath, "steamapps", "workshop", "content", "431960")
}

func GetLoginUsersPath(steamPath string) string {
	return filepath.Join(steamPath, "config", "loginusers.vdf")
}

func SteamIsRunning() bool {
	_, err := os.Stat(filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Steam", "steam.exe"))
	return err == nil
}
