package arn

// Register a list of supported anime relation types.
func init() {
	DataLists["anime-relation-types"] = []*Option{
		&Option{"prequel", HumanReadableAnimeRelation("prequel")},
		&Option{"sequel", HumanReadableAnimeRelation("sequel")},
		&Option{"alternative version", HumanReadableAnimeRelation("alternative version")},
		&Option{"alternative setting", HumanReadableAnimeRelation("alternative setting")},
		&Option{"side story", HumanReadableAnimeRelation("side story")},
		&Option{"parent story", HumanReadableAnimeRelation("parent story")},
		&Option{"full story", HumanReadableAnimeRelation("full story")},
		&Option{"spinoff", HumanReadableAnimeRelation("spinoff")},
		&Option{"summary", HumanReadableAnimeRelation("summary")},
		&Option{"other", HumanReadableAnimeRelation("other")},
	}
}

// AnimeRelation ...
type AnimeRelation struct {
	AnimeID string `json:"animeId" editable:"true"`
	Type    string `json:"type" editable:"true" datalist:"anime-relation-types"`
}

// Anime ...
func (relation *AnimeRelation) Anime() *Anime {
	anime, _ := GetAnime(relation.AnimeID)
	return anime
}

// HumanReadableType ...
func (relation *AnimeRelation) HumanReadableType() string {
	return HumanReadableAnimeRelation(relation.Type)
}

// HumanReadableAnimeRelation ...
func HumanReadableAnimeRelation(relationName string) string {
	switch relationName {
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

	return relationName
}
