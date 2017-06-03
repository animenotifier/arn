package arn

// Settings ...
type Settings struct {
	UserID        string           `json:"userId"`
	SortBy        string           `json:"sortBy"`
	TitleLanguage string           `json:"titleLanguage"`
	Providers     ServiceProviders `json:"providers"`

	user *User
}

// ServiceProviders ...
type ServiceProviders struct {
	AiringDate string `json:"airingDate"`
	Anime      string `json:"anime"`
	List       string `json:"list"`
}

// User returns the user object for the settings.
func (settings *Settings) User() *User {
	if settings.user != nil {
		return settings.user
	}

	settings.user, _ = GetUser(settings.UserID)
	return settings.user
}
