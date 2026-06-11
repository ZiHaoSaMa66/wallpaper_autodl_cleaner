package main

import (
	"fmt"
	"os"
	"wp-cleaner/steam"
)

func main() {
	steamPath, err := steam.FindSteamPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Steam path: %s\n", steamPath)
}
