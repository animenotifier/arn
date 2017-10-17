package arn

// AnimeRelations ...
type AnimeRelations struct {
	AnimeID string           `json:"animeId"`
	Items   []*AnimeRelation `json:"items"`
}
