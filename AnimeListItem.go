package arn

// AnimeListStatus values for anime list items
const (
	AnimeListStatusWatching  = "watching"
	AnimeListStatusCompleted = "completed"
	AnimeListStatusPlanned   = "planned"
	AnimeListStatusHold      = "hold"
	AnimeListStatusDropped   = "dropped"
)

// AnimeListItem ...
type AnimeListItem struct {
	AnimeID      string       `json:"animeId"`
	Status       string       `json:"status" editable:"true"`
	Episodes     int          `json:"episodes" editable:"true"`
	Rating       *AnimeRating `json:"rating"`
	Notes        string       `json:"notes" editable:"true"`
	RewatchCount int          `json:"rewatchCount" editable:"true"`
	Private      bool         `json:"private" editable:"true"`
	Created      string       `json:"created"`
	Edited       string       `json:"edited"`

	anime *Anime
}

// Anime fetches the associated anime data.
func (item *AnimeListItem) Anime() *Anime {
	if item.anime == nil {
		item.anime, _ = GetAnime(item.AnimeID)
	}

	return item.anime
}

// Link returns the URI for the given item.
func (item *AnimeListItem) Link(userNick string) string {
	return "/+" + userNick + "/animelist/" + item.AnimeID
}

// FinalRating returns the overall score for the anime.
func (item *AnimeListItem) FinalRating() float64 {
	return item.Rating.Overall
}
