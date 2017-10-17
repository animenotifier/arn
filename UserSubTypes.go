package arn

// UserAgeRange ...
type UserAgeRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// UserBrowser ...
type UserBrowser struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	IsMobile bool   `json:"isMobile"`
}

// UserOS ...
type UserOS struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// UserListProviders ...
type UserListProviders struct {
	AniList     ListProviderConfig `json:"AniList"`
	AnimePlanet ListProviderConfig `json:"AnimePlanet"`
	HummingBird ListProviderConfig `json:"HummingBird"`
	MyAnimeList ListProviderConfig `json:"MyAnimeList"`
}

// ListProviderConfig ...
type ListProviderConfig struct {
	UserName string `json:"userName"`
}

// PushEndpoint ...
type PushEndpoint struct {
	Registered string `json:"registered"`
	Keys       struct {
		P256DH string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

// UserCoverImage ...
type UserCoverImage struct {
	URL      string      `json:"url"`
	Position CSSPosition `json:"position"`
}

// CSSPosition ...
type CSSPosition struct {
	X string `json:"x"`
	Y string `json:"y"`
}

// NickToUser ...
type NickToUser struct {
	Nick   string `json:"nick"`
	UserID UserID `json:"userId"`
}

// EmailToUser ...
type EmailToUser struct {
	Email  string `json:"email"`
	UserID UserID `json:"userId"`
}

// GoogleToUser ...
type GoogleToUser struct {
	ID     string `json:"id"`
	UserID UserID `json:"userId"`
}

// FacebookToUser is the same data structure as GoogleToUser
type FacebookToUser GoogleToUser

// TwitterToUser is the same data structure as GoogleToUser
type TwitterToUser GoogleToUser
