package arn

// Genre ...
type Genre struct {
	ID        string   `json:"genre"`
	Name      string   `json:"-"`
	AnimeList []*Anime `json:"animeList"`
}

// GetGenre ...
func GetGenre(id string) (*Genre, error) {
	obj, err := DB.Get("Genre", id)
	return obj.(*Genre), err
}
