package arn

const (
	// SortByAiringDate sorts your watching list by airing date.
	SortByAiringDate = "airing date"

	// SortByTitle sorts your watching list alphabetically.
	SortByTitle = "title"

	// SortByRating sorts your watching list by rating.
	SortByRating = "rating"
)

const (
	// TitleLanguageCanonical ...
	TitleLanguageCanonical = "canonical"

	// TitleLanguageRomaji ...
	TitleLanguageRomaji = "romaji"

	// TitleLanguageEnglish ...
	TitleLanguageEnglish = "english"

	// TitleLanguageJapanese ...
	TitleLanguageJapanese = "japanese"
)

// Settings ...
type Settings struct {
	UserID        string           `json:"userId"`
	SortBy        string           `json:"sortBy"`
	TitleLanguage string           `json:"titleLanguage" editable:"true"`
	Providers     ServiceProviders `json:"providers"`

	user *User
}

// ServiceProviders ...
type ServiceProviders struct {
	Anime string `json:"anime"`
}

// NewSettings ...
func NewSettings(userID string) *Settings {
	return &Settings{
		UserID:        userID,
		SortBy:        SortByAiringDate,
		TitleLanguage: TitleLanguageCanonical,
		Providers: ServiceProviders{
			Anime: "",
		},
	}
}

// User returns the user object for the settings.
func (settings *Settings) User() *User {
	if settings.user != nil {
		return settings.user
	}

	settings.user, _ = GetUser(settings.UserID)
	return settings.user
}

// Save saves the settings in the database.
func (settings *Settings) Save() error {
	return DB.Set("Settings", settings.UserID, settings)
}
