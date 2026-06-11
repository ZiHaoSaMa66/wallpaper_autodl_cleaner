package fixer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"wp-cleaner/steam"
)

func cleanDownloads(steamPath string) error {
	dlPath := filepath.Join(steamPath, "steamapps", "workshop", "downloads")
	entries, err := os.ReadDir(dlPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("cannot read downloads dir: %w", err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "431960") {
			fullPath := filepath.Join(dlPath, e.Name())
			if err := os.RemoveAll(fullPath); err != nil {
				return fmt.Errorf("failed to remove %s: %w", e.Name(), err)
			}
			fmt.Printf("  Removed: %s\n", e.Name())
		}
	}
	return nil
}

func fixACF(steamPath string) error {
	acfPath := filepath.Join(steamPath, "steamapps", "workshop", "appworkshop_431960.acf")
	input, err := os.ReadFile(acfPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("cannot read ACF file: %w", err)
	}
	bakPath := acfPath + ".bak"
	if err := os.WriteFile(bakPath, input, 0644); err != nil {
		return fmt.Errorf("cannot create backup %s: %w", bakPath, err)
	}
	fmt.Printf("  Backup created: %s\n", bakPath)

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(string(input)))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimRight(line, "\r")
		if strings.Contains(line, "LastBuildID") {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading ACF: %w", err)
	}
	output := strings.Join(lines, "\n")
	if err := os.WriteFile(acfPath, []byte(output), 0644); err != nil {
		if renameErr := os.Rename(bakPath, acfPath); renameErr == nil {
			os.Remove(acfPath)
		}
		return fmt.Errorf("cannot write ACF: %w", err)
	}
	fmt.Println("  ACF file updated (removed LastBuildID entries)")
	return nil
}

func launchFixer(steamPath string) error {
	fixerPath := filepath.Join(steamPath, "steamapps", "common", "wallpaper_engine", "bin", "steamredownloadfixer32.exe")
	if _, err := os.Stat(fixerPath); err != nil {
		return fmt.Errorf("fixer tool not found at %s: please locate steamredownloadfixer32.exe manually and run it", fixerPath)
	}
	cmd := exec.Command(fixerPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch fixer at %s: %w: please run it manually", fixerPath, err)
	}
	fmt.Printf("Fixer launched: %s\n", fixerPath)
	fmt.Println("Please follow the on-screen instructions in the fixer tool.")
	return nil
}

func RunFixDownloads(steamPath string, force bool) error {
	fmt.Println("=== Fix Downloads: Pre-clean ===")

	running, _ := steam.CheckProcessesRunning("steam.exe", "wallpaper32.exe", "wallpaper64.exe", "ui32.exe")
	if len(running) > 0 {
		fmt.Println("The following processes are still running:")
		for _, r := range running {
			fmt.Printf("  - %s\n", r)
		}
		if !force {
			fmt.Println("Please close them before proceeding.")
			fmt.Print("Continue anyway? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Fix downloads aborted by user.")
				return nil
			}
		} else {
			fmt.Println("Continuing (force mode)...")
		}
	}

	fmt.Println("Step 1: Cleaning download cache...")
	if err := cleanDownloads(steamPath); err != nil {
		return fmt.Errorf("clean downloads failed: %w", err)
	}

	fmt.Println("Step 2: Fixing ACF file...")
	if err := fixACF(steamPath); err != nil {
		return fmt.Errorf("fix ACF failed: %w", err)
	}

	fmt.Println("Step 3: Launching Steam re-download fixer...")
	if err := launchFixer(steamPath); err != nil {
		return fmt.Errorf("launch fixer failed: %w", err)
	}

	fmt.Println("=== Fix Downloads complete ===")
	return nil
}
