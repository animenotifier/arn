package arn

// User ...
type User struct {
	ID         string       `json:"id"`
	Nick       string       `json:"nick"`
	FirstName  string       `json:"firstName"`
	LastName   string       `json:"lastName"`
	Email      string       `json:"email"`
	Role       string       `json:"role"`
	Registered string       `json:"registered"`
	LastLogin  string       `json:"lastLogin"`
	Gender     string       `json:"gender"`
	Language   string       `json:"language"`
	Avatar     string       `json:"avatar"`
	Tagline    string       `json:"tagline"`
	Website    string       `json:"website"`
	IP         string       `json:"ip"`
	UserAgent  string       `json:"agent"`
	AgeRange   UserAgeRange `json:"ageRange"`
	Location   UserLocation `json:"location"`
	Accounts   UserAccounts `json:"accounts"`
	Browser    UserBrowser  `json:"browser"`
	OS         UserOS       `json:"os"`
	Following  []string     `json:"following"`
}

// UserLocation ...
type UserLocation struct {
	CountryName string  `json:"countryName"`
	CountryCode string  `json:"countryCode"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	CityName    string  `json:"cityName"`
	RegionName  string  `json:"regionName"`
	TimeZone    string  `json:"timeZone"`
	ZipCode     string  `json:"zipCode"`
}

// UserAgeRange ...
type UserAgeRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// UserAccounts ...
type UserAccounts struct {
	Facebook struct {
		ID string `json:"id"`
	} `json:"facebook"`

	Google struct {
		ID string `json:"id"`
	} `json:"google"`

	Twitter struct {
		ID   string `json:"id"`
		Nick string `json:"nick"`
	} `json:"twitter"`

	Osu UserOsuDetails `json:"osu"`

	AniList struct {
		Nick string `json:"nick"`
	} `json:"anilist"`

	AnimePlanet struct {
		Nick string `json:"nick"`
	} `json:"animeplanet"`

	MyAnimeList struct {
		Nick string `json:"nick"`
	} `json:"myanimelist"`

	Kitsu struct {
		Nick string `json:"nick"`
	} `json:"kitsu"`
}

// UserOsuDetails ...
type UserOsuDetails struct {
	Nick     string  `json:"nick"`
	PP       float64 `json:"pp"`
	Accuracy float64 `json:"accuracy"`
	Level    float64 `json:"level"`
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

// // CoverImageStyle ...
// func (user *User) CoverImageStyle() string {
// 	url := user.CoverImage.URL

// 	if url == "" {
// 		wallpapers := []string{
// 			"https://www.pixelstalk.net/wp-content/uploads/2016/08/1080p-Anime-Desktop-Wallpaper.jpg",
// 			"https://i.imgur.com/6cJrxzx.jpg",
// 			"https://cdn.cloudpix.co/images/wallpaper-1366x768/angel-angel-beats-anime-wallpaper-666806d97b32a8a8e2b1ad9a55ab962e-large-1135606.jpg",
// 		}
// 		url = wallpapers[rand.Intn(len(wallpapers))]
// 	}

// 	return "background-image: url('" + url + "'); background-position: " + user.CoverImage.Position.X + " " + user.CoverImage.Position.Y + ";"
// }

// Settings ...
func (user *User) Settings() *Settings {
	settings := new(Settings)
	GetObject("Settings", user.ID, settings)
	return settings
}

// Save saves the user object in the database.
func (user *User) Save() {
	SetObject("User", user.ID, user)
}

// NewUser creates a new user object with default values.
func NewUser() *User {
	return &User{
	// CoverImage: UserCoverImage{
	// 	URL: "",
	// 	Position: CSSPosition{
	// 		X: "50%",
	// 		Y: "50%",
	// 	},
	// },
	}
}

// GetUser ...
func GetUser(id string) (*User, error) {
	user := NewUser()
	err := GetObject("User", id, user)
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
