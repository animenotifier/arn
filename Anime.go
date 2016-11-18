package arn

// Anime ...
type Anime struct {
	ID            int             `json:"id"`
	Type          string          `json:"type"`
	Title         AnimeTitle      `json:"title"`
	Image         string          `json:"image"`
	AiringStatus  string          `json:"airingStatus"`
	Adult         bool            `json:"adult"`
	Description   string          `json:"description"`
	StartDate     string          `json:"startDate"`
	EndDate       string          `json:"endDate"`
	Hashtag       string          `json:"hashtag"`
	YoutubeID     string          `json:"youtubeId"`
	Source        string          `json:"source"`
	TotalEpisodes int             `json:"totalEpisodes"`
	Duration      int             `json:"duration"`
	Watching      int             `json:"watching"`
	PageGenerated string          `json:"pageGenerated"`
	AnilistEdited uint64          `json:"anilistEdited"`
	Genres        []string        `json:"genres"`
	Tracks        *AnimeTrackList `json:"tracks"`
	Links         []AnimeLink     `json:"links"`
	Studios       []AnimeStudio   `json:"studios"`
	Relations     []AnimeRelation `json:"relations"`
	Created       string          `json:"created"`
	CreatedBy     string          `json:"createdBy"`
}

// AnimeTitle ...
type AnimeTitle struct {
	Romaji   string   `json:"romaji"`
	English  string   `json:"english"`
	Japanese string   `json:"japanese"`
	Synonyms []string `json:"synonyms"`
}

// AnimeLink ...
type AnimeLink struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// AnimeTrackList ...
type AnimeTrackList struct {
	Opening *AnimeTrack `json:"opening"`
}

// AnimeTrack ...
type AnimeTrack struct {
	URI        string  `json:"uri"`
	Similarity float32 `json:"similarity"`
	Permalink  string  `json:"permalink"`
	Title      string  `json:"title"`
	Likes      int     `json:"likes"`
	Plays      int     `json:"plays"`
}

// AnimeStudio ...
type AnimeStudio struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	IsMainStudio bool   `json:"isMainStudio"`
}

// AnimeRelation ...
type AnimeRelation struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

// GetAnime ...
func GetAnime(id int) (*Anime, error) {
	anime := new(Anime)
	err := GetObject("Anime", id, anime)
	return anime, err
}
