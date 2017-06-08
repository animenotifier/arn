package arn

// AnimeListStatus values for anime list items
const (
	AnimeListStatusWatching  = "watching"
	AnimeListStatusCompleted = "completed"
	AnimeListStatusPlanned   = "planned"
	AnimeListStatusDropped   = "dropped"
	AnimeListStatusHold      = "hold"
)

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

// Contains checks if the list contains the anime ID already.
func (list *AnimeList) Contains(animeID string) bool {
	for _, item := range list.Items {
		if item.AnimeID == animeID {
			return true
		}
	}

	return false
}

// Anime fetches the associated anime data.
func (item *AnimeListItem) Anime() *Anime {
	if item.anime == nil {
		item.anime, _ = GetAnime(item.AnimeID)
	}

	return item.anime
}

// Save saves the anime list in the database.
func (list *AnimeList) Save() error {
	return SetObject("AnimeList", list.UserID, list)
}

// GetAnimeList ...
func GetAnimeList(userID string) (*AnimeList, error) {
	animeList := &AnimeList{
		UserID: userID,
	}
	err := GetObject("AnimeList", userID, animeList)
	return animeList, err
}
