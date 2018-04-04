package arn

// KitsuFinder ...
type KitsuFinder struct {
	kitsuIDToAnime map[string]*Anime
}

// NewKitsuFinder creates a new finder for Kitsu anime.
func NewKitsuFinder() *KitsuFinder {
	finder := &KitsuFinder{
		kitsuIDToAnime: map[string]*Anime{},
	}

	for anime := range StreamAnime() {
		kitsuID := anime.GetMapping("kitsu/anime")

		if kitsuID != "" {
			finder.kitsuIDToAnime[kitsuID] = anime
		}
	}

	return finder
}

// GetAnime tries to find a Kitsu anime in our anime database.
func (finder *KitsuFinder) GetAnime(kitsuID string) *Anime {
	return finder.kitsuIDToAnime[kitsuID]
}
