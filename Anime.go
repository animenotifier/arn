package arn

// Anime ...
type Anime struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Title    AnimeTitle      `json:"title"`
	Image    string          `json:"image"`
	Summary  string          `json:"summary"`
	Watching int             `json:"watching"`
	Trailers []*AnimeTrailer `json:"trailers"`

	// AiringStatus  string          `json:"airingStatus"`
	// Adult         bool            `json:"adult"`
	// StartDate     string          `json:"startDate"`
	// EndDate       string          `json:"endDate"`
	// Hashtag       string          `json:"hashtag"`
	// Source        string          `json:"source"`
	// TotalEpisodes int             `json:"totalEpisodes"`
	// Duration      int             `json:"duration"`
	// PageGenerated string          `json:"pageGenerated"`
	// AnilistEdited uint64          `json:"anilistEdited"`
	// Genres        []string        `json:"genres"`
	// Tracks        *AnimeTrackList `json:"tracks"`
	// Links         []AnimeLink     `json:"links"`
	// Studios       []AnimeStudio   `json:"studios"`
	// Relations     []AnimeRelation `json:"relations"`
	// Created       string          `json:"created"`
	// CreatedBy     string          `json:"createdBy"`
}

// AnimeTitle ...
type AnimeTitle struct {
	Romaji    string   `json:"romaji"`
	English   string   `json:"english"`
	Japanese  string   `json:"japanese"`
	Canonical string   `json:"canonical"`
	Synonyms  []string `json:"synonyms"`
}

// AnimeTrailer ...
type AnimeTrailer struct {
	Service string `json:"service"`
	VideoID string `json:"videoId"`
}

// GetAnime ...
func GetAnime(id string) (*Anime, error) {
	anime := new(Anime)
	err := GetObject("Anime", id, anime)
	return anime, err
}

// Save ...
func (anime *Anime) Save() error {
	return SetObject("Anime", anime.ID, anime)
}

// FilterAnime filters all anime by a custom function.
func FilterAnime(filter func(*Anime) bool) ([]*Anime, error) {
	var filtered []*Anime

	channel := make(chan *Anime)
	err := Scan("Anime", channel)

	if err != nil {
		return filtered, err
	}

	for post := range channel {
		if filter(post) {
			filtered = append(filtered, post)
		}
	}

	return filtered, nil
}

// // GetAiringAnime ...
// func GetAiringAnime() ([]*Anime, error) {
// 	return FilterAnime(func(anime *Anime) bool {
// 		return anime.AiringStatus == "currently airing" && !anime.Adult
// 	})
// }
