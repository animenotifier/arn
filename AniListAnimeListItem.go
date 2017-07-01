package arn

import "time"

// AniListAnimeListItem ...
type AniListAnimeListItem struct {
	RecordID             int           `json:"record_id"`
	SeriesID             int           `json:"series_id"`
	ListStatus           string        `json:"list_status"`
	ScoreRaw             int           `json:"score_raw"`
	EpisodesWatched      int           `json:"episodes_watched"`
	ChaptersRead         int           `json:"chapters_read"`
	VolumesRead          int           `json:"volumes_read"`
	Rewatched            int           `json:"rewatched"`
	Reread               int           `json:"reread"`
	Priority             int           `json:"priority"`
	Private              int           `json:"private"`
	HiddenDefault        int           `json:"hidden_default"`
	Notes                string        `json:"notes"`
	AdvancedRatingScores []interface{} `json:"advanced_rating_scores"`
	CustomLists          []interface{} `json:"custom_lists"`
	StartedOn            interface{}   `json:"started_on"`
	FinishedOn           interface{}   `json:"finished_on"`
	AddedTime            time.Time     `json:"added_time"`
	UpdatedTime          time.Time     `json:"updated_time"`
	Anime                *AniListAnime `json:"anime"`
}
