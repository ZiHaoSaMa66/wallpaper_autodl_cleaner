package model

type SteamUser struct {
	SteamID64   string `json:"steamid64"`
	AccountName string `json:"account_name"`
	PersonaName string `json:"persona_name"`
	MostRecent  bool   `json:"most_recent"`
}

type WallpaperInfo struct {
	PublishedFileID uint64 `json:"publishedfileid"`
	Title           string `json:"title"`
	OwnerSteamID    string `json:"creator"`
	FileSize        uint64 `json:"file_size"`
	PreviewURL      string `json:"preview_url"`
	LocalPath       string `json:"-"`
}

type GetPublishedFileDetailsResponse struct {
	Response struct {
		PublishedFileDetails []struct {
			Result          uint64 `json:"result"`
			PublishedFileID string `json:"publishedfileid"`
			Creator         string `json:"creator"`
			CreatorAppID    uint32 `json:"creator_appid"`
			ConsumerAppID   uint32 `json:"consumer_appid"`
			Title           string `json:"title"`
			FileDescription string `json:"file_description"`
			FileSize        string `json:"file_size"`
			PreviewURL      string `json:"preview_url"`
			TimeCreated     uint32 `json:"time_created"`
			TimeUpdated     uint32 `json:"time_updated"`
		} `json:"publishedfiledetails"`
	} `json:"response"`
}

type GetUserFilesResponse struct {
	Response struct {
		Total uint64 `json:"total"`
		Files []struct {
			PublishedFileID string `json:"publishedfileid"`
		} `json:"publishedfiledetails"`
	} `json:"response"`
}
