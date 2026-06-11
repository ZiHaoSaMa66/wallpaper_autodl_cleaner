package main

import (
	"flag"
	"fmt"
	"os"
	"wp-cleaner/api"
	"wp-cleaner/cleaner"
	"wp-cleaner/model"
	"wp-cleaner/scanner"
	"wp-cleaner/steam"
)

func main() {
	apiKey := flag.String("api-key", "", "Steam Web API Key (get from https://steamcommunity.com/dev/apikey)")
	steamID := flag.String("steam-id", "", "SteamID64 to check subscriptions for (auto-detect if empty)")
	dryRun := flag.Bool("dry-run", true, "Only show what would be removed, don't actually delete")
	force := flag.Bool("force", false, "Skip confirmation prompt")
	flag.Parse()

	steamPath, err := steam.FindSteamPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Steam path: %s\n", steamPath)

	users, err := steam.GetUsers(steamPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: cannot read users: %v\n", err)
	} else {
		fmt.Println("Steam accounts on this PC:")
		for _, u := range users {
			m := ""
			if u.MostRecent {
				m = " (current)"
			}
			fmt.Printf("  %s - %s%s\n", u.AccountName, u.SteamID64, m)
		}
	}

	targetSteamID := *steamID
	if targetSteamID == "" {
		cu := steam.GetCurrentUser(users)
		if cu != nil {
			targetSteamID = cu.SteamID64
			fmt.Printf("Auto-detected current user: %s (%s)\n", cu.AccountName, cu.SteamID64)
		}
	}
	if targetSteamID == "" {
		fmt.Fprintf(os.Stderr, "ERROR: could not determine Steam user. Use -steam-id flag.\n")
		os.Exit(1)
	}

	workshopPath := steam.GetWorkshopPath(steamPath)
	localIDs, err := scanner.ScanWorkshop(workshopPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cannot scan workshop: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nLocal wallpapers found: %d\n", len(localIDs))

	fmt.Println("Fetching wallpaper metadata from Steam API...")
	var allInfos []model.WallpaperInfo
	batchSize := 100
	for i := 0; i < len(localIDs); i += batchSize {
		end := i + batchSize
		if end > len(localIDs) {
			end = len(localIDs)
		}
		infos, err := api.GetPublishedFileDetails(localIDs[i:end])
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: metadata fetch failed for batch %d: %v\n", i/batchSize, err)
			continue
		}
		allInfos = append(allInfos, infos...)
	}
	fmt.Printf("Metadata retrieved: %d\n", len(allInfos))

	fmt.Println("Fetching subscription list...")
	subscribedIDs, err := steam.GetSubscribedIDs(*apiKey, targetSteamID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: %v\n", err)
		fmt.Println("Falling back: all wallpapers will be listed. Subscribe list unavailable.")
		subscribedIDs = nil
	}
	if subscribedIDs != nil {
		fmt.Printf("Subscribed wallpapers: %d\n", len(subscribedIDs))
	}

	plan := cleaner.BuildCleanupPlan(workshopPath, localIDs, subscribedIDs, allInfos)
	plan.DryRun()

	if !*dryRun {
		if !*force {
			fmt.Print("\nProceed with cleanup? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cancelled.")
				return
			}
		}
		if err := plan.Execute(false); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR during cleanup: %v\n", err)
			os.Exit(1)
		}
	}
}
