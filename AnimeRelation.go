package arn

// AnimeRelation ...
type AnimeRelation struct {
	AnimeID string `json:"animeId"`
	Type    string `json:"type"`

	anime *Anime
}

// Anime ...
func (relation *AnimeRelation) Anime() *Anime {
	if relation.anime == nil {
		relation.anime, _ = GetAnime(relation.AnimeID)
	}

	return relation.anime
}

// HumanReadableType ...
func (relation *AnimeRelation) HumanReadableType() string {
	switch relation.Type {
	case "prequel":
		return "Prequel"
	case "sequel":
		return "Sequel"
	case "alternative version":
		return "Alternative"
	case "alternative setting":
		return "Alternative"
	case "side story":
		return "Side story"
	case "parent story":
		return "Parent story"
	case "full story":
		return "Full story"
	case "spinoff":
		return "Spin-off"
	case "summary":
		return "Summary"
	case "other":
		return "Other"
	}

	return relation.Type
}
