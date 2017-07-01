package arn

// AniListAnime ...
type AniListAnime struct {
	ID             int           `json:"id"`
	TitleRomaji    string        `json:"title_romaji"`
	TitleEnglish   string        `json:"title_english"`
	TitleJapanese  string        `json:"title_japanese"`
	Type           string        `json:"type"`
	StartDateFuzzy int           `json:"start_date_fuzzy"`
	EndDateFuzzy   int           `json:"end_date_fuzzy"`
	Season         interface{}   `json:"season"`
	SeriesType     string        `json:"series_type"`
	Synonyms       []interface{} `json:"synonyms"`
	Genres         []string      `json:"genres"`
	Adult          bool          `json:"adult"`
	AverageScore   float64       `json:"average_score"`
	Popularity     int           `json:"popularity"`
	UpdatedAt      int           `json:"updated_at"`
	Hashtag        interface{}   `json:"hashtag"`
	ImageURLSml    string        `json:"image_url_sml"`
	ImageURLMed    string        `json:"image_url_med"`
	ImageURLLge    string        `json:"image_url_lge"`
	ImageURLBanner interface{}   `json:"image_url_banner"`
	TotalEpisodes  int           `json:"total_episodes"`
	AiringStatus   string        `json:"airing_status"`
}
