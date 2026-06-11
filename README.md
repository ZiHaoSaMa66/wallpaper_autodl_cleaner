# wp-cleaner

> Clean up wallpapers auto-downloaded by other Steam accounts in Wallpaper Engine

English [中文](README_ZH.md)

## Usage

### 1. Get a Steam Web API Key
Go to https://steamcommunity.com/dev/apikey and register a key (free, any domain name works).

### 2. Run the tool
```cmd
wp-cleaner.exe -api-key=YOUR_STEAM_API_KEY -dry-run
```

This will:
- Detect your Steam installation path
- Identify the current logged-in Steam user
- Scan all downloaded wallpapers in workshop/content/431960/
- Fetch wallpaper metadata (titles, owners) from Steam
- Compare against your subscription list
- Show which wallpapers can be safely removed

### 3. Execute cleanup
```cmd
wp-cleaner.exe -api-key=YOUR_STEAM_API_KEY
```

Removes (moves to hidden backup) all wallpapers not subscribed by the current user.

### 4. Fix download issues (optional)
```cmd
wp-cleaner.exe -fix-downloads
```

Clears Steam download cache, fixes ACF files, and launches the Steam re-download fixer tool. Run this before cleanup if you have download issues.

### 5. Permanently delete quarantined trash
```cmd
wp-cleaner.exe -api-key=YOUR_STEAM_API_KEY -delete
```

After cleanup, quarantined `.trash-*` folders remain for review. Add `-delete` to permanently remove them at the end of cleanup.

### 6. Full workflow
```cmd
wp-cleaner.exe -api-key=YOUR_STEAM_API_KEY -fix-downloads -delete -dry-run=false
```

Fix downloads → clean up unsubscribed wallpapers → permanently delete trash.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-api-key` | `""` | Steam Web API Key (required for subscription detection) |
| `-steam-id` | `""` | SteamID64 (auto-detected from loginusers.vdf if empty) |
| `-dry-run` | `true` | Preview only, no actual deletion |
| `-force` | `false` | Skip confirmation prompt |
| `-fix-downloads` | `false` | Pre-cleanup: clear download cache, fix ACF files, run Steam repair tool |
| `-delete` | `false` | Permanently delete quarantined .trash-* folders (runs normal cleanup first) |

## How it works

1. Reads Steam installation path from Windows registry
2. Parses `loginusers.vdf` to find all Steam accounts on this PC
3. Scans `steamapps/workshop/content/431960/` for downloaded wallpaper folders
4. Calls `IPublishedFileService/GetUserFiles?type=mysubscriptions` to get YOUR subscribed items
5. Calls `GetPublishedFileDetails` (public API) to get wallpaper titles
6. Compares lists and identifies wallpapers NOT in your subscriptions
7. Renames non-subscribed folders to `.trash-*` for safe review/deletion

## Build from source

```bash
go build -o wp-cleaner.exe .
```

Requires: Go 1.22+, Windows 10/11
