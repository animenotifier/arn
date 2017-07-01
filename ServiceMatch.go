package arn

// ServiceMatch ...
type ServiceMatch struct {
	AnimeID    string  `json:"animeId"`
	ServiceID  string  `json:"providerId"`
	Similarity float64 `json:"similarity"`
	Edited     string  `json:"edited"`
	EditedBy   string  `json:"editedBy"`
}

// AniListToAnime ...
type AniListToAnime ServiceMatch

// GetAniListToAnime ...
func GetAniListToAnime(aniListID string) (*AniListToAnime, error) {
	obj, err := DB.Get("AniListToAnime", aniListID)
	return obj.(*AniListToAnime), err
}
