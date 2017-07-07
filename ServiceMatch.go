package arn

// ServiceMatch ...
type ServiceMatch struct {
	AnimeID    string  `json:"animeId"`
	ServiceID  string  `json:"serviceId"`
	Similarity float64 `json:"similarity"`
	Edited     string  `json:"edited"`
	EditedBy   string  `json:"editedBy"`
}

// AniListToAnime ...
type AniListToAnime ServiceMatch

// MyAnimeListToAnime ...
type MyAnimeListToAnime ServiceMatch

// GetAniListToAnime ...
func GetAniListToAnime(aniListID string) (*AniListToAnime, error) {
	obj, err := DB.Get("AniListToAnime", aniListID)
	return obj.(*AniListToAnime), err
}

// GetMyAnimeListToAnime ...
func GetMyAnimeListToAnime(malID string) (*MyAnimeListToAnime, error) {
	obj, err := DB.Get("MyAnimeListToAnime", malID)
	return obj.(*MyAnimeListToAnime), err
}
