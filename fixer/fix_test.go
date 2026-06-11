package fixer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanDownloads_Removes431960(t *testing.T) {
	steamPath := t.TempDir()
	dlPath := filepath.Join(steamPath, "steamapps", "workshop", "downloads")
	os.MkdirAll(dlPath, 0755)

	os.MkdirAll(filepath.Join(dlPath, "431960_12345"), 0755)
	os.MkdirAll(filepath.Join(dlPath, "other_app"), 0755)

	err := cleanDownloads(steamPath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dlPath, "431960_12345")); !os.IsNotExist(err) {
		t.Error("431960_12345 should have been removed")
	}
	if _, err := os.Stat(filepath.Join(dlPath, "other_app")); os.IsNotExist(err) {
		t.Error("other_app should still exist")
	}
}

func TestCleanDownloads_NoDir(t *testing.T) {
	steamPath := t.TempDir()
	err := cleanDownloads(steamPath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFixACF_RemovesLastBuildID(t *testing.T) {
	steamPath := t.TempDir()
	acfDir := filepath.Join(steamPath, "steamapps", "workshop")
	os.MkdirAll(acfDir, 0755)
	acfPath := filepath.Join(acfDir, "appworkshop_431960.acf")

	input := `"AppWorkshop"
{
	"appid"		"431960"
	"LastBuildID"	"12345"
	"SizeOnDisk"	"1000"
	"LastBuildID"	"67890"
}
`
	if err := os.WriteFile(acfPath, []byte(input), 0644); err != nil {
		t.Fatal(err)
	}

	err := fixACF(steamPath)
	if err != nil {
		t.Fatal(err)
	}

	output, err := os.ReadFile(acfPath)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(output), "LastBuildID") {
		t.Error("LastBuildID should have been removed from ACF file")
	}
	if !strings.Contains(string(output), "SizeOnDisk") {
		t.Error("non-LastBuildID entries should remain")
	}

	if _, err := os.Stat(acfPath + ".bak"); os.IsNotExist(err) {
		t.Error("backup file should exist")
	}
}

func TestFixACF_NoFile(t *testing.T) {
	steamPath := t.TempDir()
	err := fixACF(steamPath)
	if err != nil {
		t.Fatal(err)
	}
}
