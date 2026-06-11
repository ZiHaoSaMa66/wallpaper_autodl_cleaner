package main

import (
	"testing"
	"wp-cleaner/steam"
)

func TestWorkshopPath_UsesWELibrary(t *testing.T) {
	weLibraryPath := `E:\SteamLibrary`
	want := `E:\SteamLibrary\steamapps\workshop\content\431960`
	got := steam.GetWorkshopPath(weLibraryPath)
	if got != want {
		t.Fatalf("GetWorkshopPath(%q) = %q, want %q", weLibraryPath, got, want)
	}
}
