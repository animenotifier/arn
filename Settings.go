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
	UserID        UserID           `json:"userId"`
	SortBy        string           `json:"sortBy"`
	TitleLanguage string           `json:"titleLanguage" editable:"true"`
	Providers     ServiceProviders `json:"providers"`
	Avatar        AvatarSettings   `json:"avatar"`
	Format        FormatSettings   `json:"format"`

	user *User
}

// FormatSettings ...
type FormatSettings struct {
	RatingsPrecision int `json:"ratingsPrecision" editable:"true"`
}

// ServiceProviders ...
type ServiceProviders struct {
	Anime string `json:"anime"`
}

// AvatarSettings ...
type AvatarSettings struct {
	Source    string `json:"source" editable:"true"`
	SourceURL string `json:"sourceUrl" editable:"true"`
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
		Avatar: AvatarSettings{
			Source:    "",
			SourceURL: "",
		},
	}
}

// GetSettings ...
func GetSettings(userID string) (*Settings, error) {
	obj, err := DB.Get("Settings", userID)

	if err != nil {
		return nil, err
	}

	return obj.(*Settings), nil
}

// User returns the user object for the settings.
func (settings *Settings) User() *User {
	if settings.user != nil {
		return settings.user
	}

	settings.user, _ = GetUser(settings.UserID)
	return settings.user
}
