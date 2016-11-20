package arn

// Genre ...
type Genre struct {
	ID        string   `json:"genre"`
	Name      string   `json:"-"`
	AnimeList []*Anime `json:"animeList"`
}

// GetGenre ...
func GetGenre(id string) (*Genre, error) {
	genre := new(Genre)
	err := GetObject("Genres", id, genre)
	return genre, err
}
