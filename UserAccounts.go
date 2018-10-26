package arn

// UserAccounts ...
type UserAccounts struct {
	Facebook struct {
		ID string `json:"id" private:"true"`
	} `json:"facebook"`

	Google struct {
		ID string `json:"id" private:"true"`
	} `json:"google"`

	Twitter struct {
		ID   string `json:"id" private:"true"`
		Nick string `json:"nick"`
	} `json:"twitter"`

	Osu struct {
		Nick     string  `json:"nick" editable:"true"`
		PP       float64 `json:"pp"`
		Accuracy float64 `json:"accuracy"`
		Level    float64 `json:"level"`
	} `json:"osu"`

	Overwatch struct {
		BattleTag   string `json:"battleTag" editable:"true"`
		SkillRating int    `json:"skillRating"`
		Tier        string `json:"tier"`
	} `json:"overwatch"`

	FinalFantasyXIV struct {
		Name      string `json:"name" editable:"true"`
		Server    string `json:"server" editable:"true"`
		Class     string `json:"class"`
		ItemLevel int    `json:"itemLevel"`
	} `json:"ffxiv"`

	AniList struct {
		Nick string `json:"nick" editable:"true"`
	} `json:"anilist"`

	AnimePlanet struct {
		Nick string `json:"nick" editable:"true"`
	} `json:"animeplanet"`

	MyAnimeList struct {
		Nick string `json:"nick" editable:"true"`
	} `json:"myanimelist"`

	Kitsu struct {
		Nick string `json:"nick" editable:"true"`
	} `json:"kitsu"`
}
