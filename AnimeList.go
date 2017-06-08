package arn

// AnimeList ...
type AnimeList struct {
	UserID string          `json:"userId"`
	Items  []AnimeListItem `json:"items"`
}

// AnimeListItem ...
type AnimeListItem struct {
	AnimeID      string      `json:"animeId"`
	Status       string      `json:"status"`
	Episode      int         `json:"episode"`
	Rating       AnimeRating `json:"rating"`
	Notes        string      `json:"notes"`
	RewatchCount int         `json:"rewatchCount"`
	Private      bool        `json:"private"`

	anime *Anime
}

// Anime fetches the associated anime data.
func (item *AnimeListItem) Anime() *Anime {
	if item.anime == nil {
		item.anime, _ = GetAnime(item.AnimeID)
	}

	return item.anime
}

// GetAnimeList ...
func GetAnimeList(userID string) (*AnimeList, error) {
	animeList := new(AnimeList)
	err := GetObject("AnimeList", userID, animeList)
	return animeList, err
}
