package arn

// AMV is an anime music video.
type AMV struct {
	ID          string   `json:"id"`
	Title       AMVTitle `json:"title"`
	MainAnimeID string   `json:"mainAnimeId" editable:"true"`
	AnimeIDs    []string `json:"animeIds" editable:"true"`
	Tags        []string `json:"tags" editable:"true"`
	IsDraft     bool     `json:"isDraft" editable:"true"`

	HasCreator
	HasEditor
	HasLikes
}
