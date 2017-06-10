package arn

import (
	"sync"
)

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

	animeList := make([]*Anime, len(cache.IDList))

	var wg sync.WaitGroup
	wg.Add(len(cache.IDList))

	for index, id := range cache.IDList {
		listIndex := index
		animeID := id

		go func() {
			anime, getErr := GetAnime(animeID)

			if anime == nil || getErr != nil {
				animeList[listIndex] = NotFoundAnime
			} else {
				animeList[listIndex] = anime
			}

			wg.Done()
		}()
	}

	wg.Wait()

	return animeList, nil
}
