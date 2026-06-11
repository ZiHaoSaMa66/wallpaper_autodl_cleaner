package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"wp-cleaner/model"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

const (
	PublicAPI = "https://api.steampowered.com"
)

func GetPublishedFileDetails(ids []uint64) ([]model.WallpaperInfo, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	if len(ids) > 100 {
		ids = ids[:100]
	}

	form := url.Values{}
	form.Set("itemcount", fmt.Sprintf("%d", len(ids)))
	for i, id := range ids {
		form.Set(fmt.Sprintf("publishedfileids[%d]", i), fmt.Sprintf("%d", id))
	}

	resp, err := httpClient.PostForm(
		PublicAPI+"/ISteamRemoteStorage/GetPublishedFileDetails/v1/",
		form,
	)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	var apiResp model.GetPublishedFileDetailsResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	var results []model.WallpaperInfo
	for _, d := range apiResp.Response.PublishedFileDetails {
		if d.Result != 1 {
			continue
		}
		id, err := parseUint64(d.PublishedFileID)
		if err != nil {
			continue
		}
		var fileSize uint64
		if d.FileSize != "" {
			fileSize, _ = strconv.ParseUint(d.FileSize, 10, 64)
		}
		results = append(results, model.WallpaperInfo{
			PublishedFileID: id,
			Title:           d.Title,
			OwnerSteamID:    d.Creator,
			FileSize:        fileSize,
			PreviewURL:      d.PreviewURL,
		})
	}
	return results, nil
}

func GetUserSubscribedFiles(apiKey, steamID string) ([]uint64, error) {
	inputJSON := fmt.Sprintf(
		`{"steamid":"%s","appid":%s,"type":"mysubscriptions","numperpage":100}`,
		steamID, model.WallpaperEngineAppID,
	)
	reqURL := fmt.Sprintf(
		"%s/IPublishedFileService/GetUserFiles/v1/?key=%s&input_json=%s",
		PublicAPI, apiKey, url.QueryEscape(inputJSON),
	)

	resp, err := httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("subscription API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read subscription response failed: %w", err)
	}

	var apiResp model.GetUserFilesResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse subscription response failed: %w", err)
	}

	var ids []uint64
	for _, f := range apiResp.Response.Files {
		id, err := parseUint64(f.PublishedFileID)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	total := apiResp.Response.Total
	if total > 100 {
		page := uint64(2)
		for uint64(len(ids)) < total {
			inputJSON := fmt.Sprintf(
				`{"steamid":"%s","appid":%s,"type":"mysubscriptions","numperpage":100,"page":%d}`,
				steamID, model.WallpaperEngineAppID, page,
			)
			reqURL := fmt.Sprintf(
				"%s/IPublishedFileService/GetUserFiles/v1/?key=%s&input_json=%s",
				PublicAPI, apiKey, url.QueryEscape(inputJSON),
			)
			resp, err := httpClient.Get(reqURL)
			if err != nil {
				return ids, fmt.Errorf("pagination page %d failed: %w", page, err)
			}
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return ids, fmt.Errorf("pagination page %d read failed: %w", page, err)
			}
			var pageResp model.GetUserFilesResponse
			if err := json.Unmarshal(body, &pageResp); err != nil {
				return ids, fmt.Errorf("pagination page %d parse failed: %w", page, err)
			}
			for _, f := range pageResp.Response.Files {
				id, err := parseUint64(f.PublishedFileID)
				if err != nil {
					continue
				}
				ids = append(ids, id)
			}
			page++
		}
	}

	return ids, nil
}

func parseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
