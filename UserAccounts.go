package arn

// UserAccounts ...
type UserAccounts struct {
	Facebook struct {
		ID string `json:"id"`
	} `json:"facebook"`

	Google struct {
		ID string `json:"id"`
	} `json:"google"`

	Twitter struct {
		ID   string `json:"id"`
		Nick string `json:"nick"`
	} `json:"twitter"`

	Osu UserOsuDetails `json:"osu"`

	AniList struct {
		Nick string `json:"nick"`
	} `json:"anilist"`

	AnimePlanet struct {
		Nick string `json:"nick"`
	} `json:"animeplanet"`

	MyAnimeList struct {
		Nick string `json:"nick"`
	} `json:"myanimelist"`

	Kitsu struct {
		Nick string `json:"nick"`
	} `json:"kitsu"`
}
