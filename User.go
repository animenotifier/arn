package arn

import "math/rand"

// User ...
type User struct {
	ID   string `json:"id"`
	Nick string `json:"nick"`
	// 	FirstName     string            `json:"firstName"`
	// 	LastName      string            `json:"lastName"`
	// 	Email         string            `json:"email"`
	// 	Gender        string            `json:"gender"`
	// 	Language      string            `json:"language"`
	Avatar string `json:"avatar"`
	// 	Osu           string            `json:"osu"`
	Registered string `json:"registered"`
	Role       string `json:"role"`
	// 	SortBy        string            `json:"sortBy"`
	Tagline string `json:"tagline"`
	// 	TitleLanguage string            `json:"titleLanguage"`
	// 	Twitter       string            `json:"twitter"`
	// 	Website       string            `json:"website"`
	// 	IP            string            `json:"ip"`
	// 	LastLogin     string            `json:"lastLogin"`
	// 	Providers     UserProviders     `json:"providers"`
	// 	ListProviders UserListProviders `json:"listProviders"`
	// 	Accounts      UserAccounts      `json:"accounts"`
	// 	AgeRange      struct {
	// 		Min int `json:"min"`
	// 		Max int `json:"max"`
	// 	} `json:"ageRange"`
	CoverImage UserCoverImage `json:"coverImage"`
	// 	Agent         UserAgent               `json:"agent"`
	// 	Location      UserLocation            `json:"location"`
	// 	OsuDetails    UserOsuDetails          `json:"osuDetails"`
	// 	Following     []string                `json:"following"`
	// 	PushEndpoints map[string]PushEndpoint `json:"pushEndpoints"`
	// 	LastView      struct {
	// 		Date string `json:"date"`
	// 		URL  string `json:"url"`
	// 	} `json:"lastView"`
	// }

	// // UserLocation ...
	// type UserLocation struct {
	// 	CountryName   string `json:"countryName"`
	// 	CountryCode   string `json:"countryCode"`
	// 	Latitude      string `json:"latitude"`
	// 	Longitude     string `json:"longitude"`
	// 	IPAddress     string `json:"ipAddress"`
	// 	ZipCode       string `json:"zipCode"`
	// 	CityName      string `json:"cityName"`
	// 	TimeZone      string `json:"timeZone"`
	// 	RegionName    string `json:"regionName"`
	// 	StatusCode    string `json:"statusCode"`
	// 	StatusMessage string `json:"statusMessage"`
}

// UserAccounts ...
type UserAccounts struct {
	Facebook string `json:"facebook"`
	Google   string `json:"google"`
	Twitter  int    `json:"twitter"`
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

// CoverImageStyle ...
func (user *User) CoverImageStyle() string {
	url := user.CoverImage.URL

	if url == "" {
		wallpapers := []string{
			"https://www.pixelstalk.net/wp-content/uploads/2016/08/1080p-Anime-Desktop-Wallpaper.jpg",
			"https://i.imgur.com/6cJrxzx.jpg",
			"https://cdn.cloudpix.co/images/wallpaper-1366x768/angel-angel-beats-anime-wallpaper-666806d97b32a8a8e2b1ad9a55ab962e-large-1135606.jpg",
		}
		url = wallpapers[rand.Intn(len(wallpapers))]
	}

	return "background-image: url('" + url + "'); background-position: " + user.CoverImage.Position.X + " " + user.CoverImage.Position.Y + ";"
}

// Save saves the user object in the database.
func (user *User) Save() {
	SetObject("Users", user.ID, user)
}

// NewUser creates a new user object with default values.
func NewUser() *User {
	return &User{
		CoverImage: UserCoverImage{
			URL: "",
			Position: CSSPosition{
				X: "50%",
				Y: "50%",
			},
		},
	}
}

// GetUser ...
func GetUser(id string) (*User, error) {
	user := NewUser()
	err := GetObject("Users", id, user)
	return user, err
}

// GetUserByNick ...
func GetUserByNick(nick string) (*User, error) {
	rec, err := Get("NickToUser", nick)

	if err != nil {
		return nil, err
	}

	return GetUser(rec["userId"].(string))
}
