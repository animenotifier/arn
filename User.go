package arn

// User ...
type User struct {
	ID            string            `json:"id"`
	Nick          string            `json:"nick"`
	FirstName     string            `json:"firstName"`
	LastName      string            `json:"lastName"`
	Email         string            `json:"email"`
	Gender        string            `json:"gender"`
	Language      string            `json:"language"`
	Avatar        string            `json:"avatar"`
	Osu           string            `json:"osu"`
	Registered    string            `json:"registered"`
	Role          string            `json:"role"`
	SortBy        string            `json:"sortBy"`
	Tagline       string            `json:"tagline"`
	TitleLanguage string            `json:"titleLanguage"`
	Twitter       string            `json:"twitter"`
	Website       string            `json:"website"`
	IP            string            `json:"ip"`
	LastLogin     string            `json:"lastLogin"`
	Providers     UserProviders     `json:"providers"`
	ListProviders UserListProviders `json:"listProviders"`
	Accounts      UserAccounts      `json:"accounts"`
	AgeRange      struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"ageRange"`
	Agent         UserAgent               `json:"agent"`
	Location      UserLocation            `json:"location"`
	OsuDetails    UserOsuDetails          `json:"osuDetails"`
	Following     []string                `json:"following"`
	PushEndpoints map[string]PushEndpoint `json:"pushEndpoints"`
	LastView      struct {
		Date string `json:"date"`
		URL  string `json:"url"`
	} `json:"lastView"`
}

// UserLocation ...
type UserLocation struct {
	CountryName   string `json:"countryName"`
	CountryCode   string `json:"countryCode"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	IPAddress     string `json:"ipAddress"`
	ZipCode       string `json:"zipCode"`
	CityName      string `json:"cityName"`
	TimeZone      string `json:"timeZone"`
	RegionName    string `json:"regionName"`
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

// UserAccounts ...
type UserAccounts struct {
	Facebook string `json:"facebook"`
	Google   string `json:"google"`
	Twitter  string `json:"twitter"`
}

// UserAgent ...
type UserAgent struct {
	Family string `json:"family"`
	Patch  string `json:"patch"`
	Minor  string `json:"minor"`
	Major  string `json:"major"`
	Source string `json:"source"`
}

// UserOsuDetails ...
type UserOsuDetails struct {
	PP        float64 `json:"pp"`
	Level     float64 `json:"level"`
	Nick      string  `json:"nick"`
	Accuracy  float64 `json:"accuracy"`
	PlayCount int     `json:"playCount"`
}

// UserProviders ...
type UserProviders struct {
	AiringDate string `json:"airingDate"`
	Anime      string `json:"anime"`
	List       string `json:"list"`
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
