package arn

// AnimeRelations ...
type AnimeRelations struct {
	AnimeID AnimeID          `json:"animeId"`
	Items   []*AnimeRelation `json:"items"`
}
