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

	if err != nil {
		return nil, err
	}

	return obj.(*AniListToAnime), nil
}

// GetMyAnimeListToAnime ...
func GetMyAnimeListToAnime(malID string) (*MyAnimeListToAnime, error) {
	obj, err := DB.Get("MyAnimeListToAnime", malID)

	if err != nil {
		return nil, err
	}

	return obj.(*MyAnimeListToAnime), nil
}
