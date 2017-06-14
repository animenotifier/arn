package arn

// ListOfIDs ...
type ListOfIDs struct {
	IDList []string `json:"idList"`
}

// GetAiringAnimeCached ...
func GetAiringAnimeCached() ([]*Anime, error) {
	var cache ListOfIDs

	err := DB.GetObject("Cache", "airing anime", &cache)

	if err != nil {
		return nil, err
	}

	animeList, err := DB.GetMany("Anime", cache.IDList)
	return animeList.([]*Anime), err
}
