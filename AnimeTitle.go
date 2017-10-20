package arn

// AnimeTitle ...
type AnimeTitle struct {
	Romaji    string   `json:"romaji"`
	English   string   `json:"english"`
	Japanese  string   `json:"japanese"`
	Hiragana  string   `json:"hiragana"`
	Canonical string   `json:"canonical"`
	Synonyms  []string `json:"synonyms"`
}

// ByUser ...
func (title *AnimeTitle) ByUser(user *User) string {
	if user == nil {
		return title.Canonical
	}

	switch user.Settings().TitleLanguage {
	case "canonical":
		return title.Canonical
	case "romaji":
		if title.Romaji == "" {
			return title.Canonical
		}

		return title.Romaji
	case "english":
		if title.English == "" {
			return title.Canonical
		}

		return title.English
	case "japanese":
		if title.Japanese == "" {
			return title.Canonical
		}

		return title.Japanese
	default:
		panic("Invalid title language")
	}
}
