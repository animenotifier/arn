package arn

// DefaultAverageRating is the average rating we're going to assume for an anime with 0 ratings.
const DefaultAverageRating = 5.0

// AnimeRating ...
type AnimeRating struct {
	Overall    float64 `json:"overall" editable:"true"`
	Story      float64 `json:"story" editable:"true"`
	Visuals    float64 `json:"visuals" editable:"true"`
	Soundtrack float64 `json:"soundtrack" editable:"true"`
}

// IsNotRated tells you whether all ratings are zero.
func (rating *AnimeRating) IsNotRated() bool {
	return rating.Overall == 0 && rating.Story == 0 && rating.Visuals == 0 && rating.Soundtrack == 0
}

// Reset sets all values to the default anime average rating.
func (rating *AnimeRating) Reset() {
	rating.Overall = DefaultAverageRating
	rating.Story = DefaultAverageRating
	rating.Visuals = DefaultAverageRating
	rating.Soundtrack = DefaultAverageRating
}
