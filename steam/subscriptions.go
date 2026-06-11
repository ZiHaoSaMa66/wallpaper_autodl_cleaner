package steam

import (
	"fmt"
	"wp-cleaner/api"
)

func GetSubscribedIDs(apiKey, steamID string) ([]uint64, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("no API key provided; run with -api-key to auto-detect subscriptions")
	}
	ids, err := api.GetUserSubscribedFiles(apiKey, steamID)
	if err != nil {
		return nil, fmt.Errorf("API subscription query failed: %w\nTip: get a free Steam Web API key at https://steamcommunity.com/dev/apikey", err)
	}
	return ids, nil
}

func BuildSubscribedSet(ids []uint64) map[uint64]bool {
	set := make(map[uint64]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	return set
}
