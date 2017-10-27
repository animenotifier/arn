package arn

// AnimeRelations ...
type AnimeRelations struct {
	AnimeID string           `json:"animeId"`
	Items   []*AnimeRelation `json:"items"`
}

// GetAnimeRelations ...
func GetAnimeRelations(animeID string) (*AnimeRelations, error) {
	obj, err := DB.Get("AnimeRelations", animeID)

	if err != nil {
		return nil, err
	}

	return obj.(*AnimeRelations), nil
}
