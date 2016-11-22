package arn

// PopularCache ...
type PopularCache struct {
	Anime []*Anime `json:"anime"`
}

// GetPopularCache ...
func GetPopularCache() (*PopularCache, error) {
	cache := new(PopularCache)
	err := GetObject("Cache", "popularAnime", cache)
	return cache, err
}
