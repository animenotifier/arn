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
	switch user.Settings().TitleLanguage {
	case "canonical":
		return title.Canonical
	case "romaji":
		return title.Romaji
	case "english":
		if title.English == "" {
			return title.Canonical
		}

		return title.English
	case "japanese":
		return title.Japanese
	default:
		panic("Invalid title language")
	}
}
