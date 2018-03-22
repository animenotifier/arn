package arn

// DefaultRating is the default rating value.
const DefaultRating = 0.0

// AverageRating is the center rating in the system.
// Note that the mathematically correct center would be a little higher,
// but we don't care about these slight offsets.
const AverageRating = 5.0

// MaxRating is the maximum rating users can give.
const MaxRating = 10.0

// RatingCountThreshold is the number of users threshold that, when passed, doesn't dampen the result.
const RatingCountThreshold = 4

// AnimeRating ...
type AnimeRating struct {
	Overall    float64 `json:"overall" editable:"true"`
	Story      float64 `json:"story" editable:"true"`
	Visuals    float64 `json:"visuals" editable:"true"`
	Soundtrack float64 `json:"soundtrack" editable:"true"`

	// The amount of people who rated
	Count AnimeRatingCount `json:"count"`
}

// AnimeRatingCount ...
type AnimeRatingCount struct {
	Overall    int `json:"overall"`
	Story      int `json:"story"`
	Visuals    int `json:"visuals"`
	Soundtrack int `json:"soundtrack"`
}

// IsNotRated tells you whether all ratings are zero.
func (rating *AnimeRating) IsNotRated() bool {
	return rating.Overall == 0 && rating.Story == 0 && rating.Visuals == 0 && rating.Soundtrack == 0
}

// Reset sets all values to the default anime average rating.
func (rating *AnimeRating) Reset() {
	rating.Overall = DefaultRating
	rating.Story = DefaultRating
	rating.Visuals = DefaultRating
	rating.Soundtrack = DefaultRating
}

// Clamp ...
func (rating *AnimeRating) Clamp() {
	if rating.Overall < 0 {
		rating.Overall = 0
	}

	if rating.Story < 0 {
		rating.Story = 0
	}

	if rating.Visuals < 0 {
		rating.Visuals = 0
	}

	if rating.Soundtrack < 0 {
		rating.Soundtrack = 0
	}

	if rating.Overall > MaxRating {
		rating.Overall = MaxRating
	}

	if rating.Story > MaxRating {
		rating.Story = MaxRating
	}

	if rating.Visuals > MaxRating {
		rating.Visuals = MaxRating
	}

	if rating.Soundtrack > MaxRating {
		rating.Soundtrack = MaxRating
	}
}
