package steam

import (
	"fmt"
	"os"

	"github.com/jslay88/vdf"
	"wp-cleaner/model"
)

func GetUsers(steamPath string) ([]model.SteamUser, error) {
	path := GetLoginUsersPath(steamPath)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	doc, err := vdf.ParseString(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing loginusers.vdf: %w", err)
	}

	usersNode := doc.Get("users")
	if usersNode == nil {
		return nil, nil
	}

	var users []model.SteamUser
	for _, child := range usersNode.Children {
		if !child.IsObject {
			continue
		}
		u := model.SteamUser{SteamID64: child.Key}
		u.AccountName = child.GetString("AccountName")
		u.PersonaName = child.GetString("PersonaName")
		if child.GetString("MostRecent") == "1" {
			u.MostRecent = true
		}
		users = append(users, u)
	}
	return users, nil
}

func GetCurrentUser(users []model.SteamUser) *model.SteamUser {
	for i, u := range users {
		if u.MostRecent {
			return &users[i]
		}
	}
	if len(users) > 0 {
		return &users[0]
	}
	return nil
}
