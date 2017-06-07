package arn

import (
	"time"
)

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

// NickToUser ...
type NickToUser struct {
	Nick   string `json:"nick"`
	UserID string `json:"userId"`
}

// EmailToUser ...
type EmailToUser struct {
	Email  string `json:"email"`
	UserID string `json:"userId"`
}

// CoverImageURL ...
func (user *User) CoverImageURL() string {
	return "/images/cover/default"
	// url := user.CoverImage.URL
	// url := ""

	// if url == "" {
	// 	wallpapers := []string{
	// 		"https://www.pixelstalk.net/wp-content/uploads/2016/08/1080p-Anime-Desktop-Wallpaper.jpg",
	// 		"https://i.imgur.com/6cJrxzx.jpg",
	// 		"https://cdn.cloudpix.co/images/wallpaper-1366x768/angel-angel-beats-anime-wallpaper-666806d97b32a8a8e2b1ad9a55ab962e-large-1135606.jpg",
	// 		"https://s-media-cache-ak0.pinimg.com/originals/26/bc/e8/26bce85b5a225f294859ff9be7ba7326.jpg",
	// 	}
	// 	url = wallpapers[rand.Intn(len(wallpapers))]
	// }

	// return url
	//  background-position: " + user.CoverImage.Position.X + " " + user.CoverImage.Position.Y + ";"
}

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

// SetNick changes the user's nickname safely.
func (user *User) SetNick(newName string) {
	// Delete old nick reference
	Delete("NickToUser", user.Nick)

	// Set new nick
	user.Nick = newName

	// New nick reference
	record := &NickToUser{
		Nick:   user.Nick,
		UserID: user.ID,
	}

	SetObject("NickToUser", record.Nick, record)
}

// SetEmail changes the user's email safely.
func (user *User) SetEmail(newName string) {
	// Delete old email reference
	Delete("EmailToUser", user.Email)

	// Set new email
	user.Email = newName

	// New email reference
	record := &EmailToUser{
		Email:  user.Email,
		UserID: user.ID,
	}

	SetObject("EmailToUser", record.Email, record)
}

// RegisteredTime ...
func (user *User) RegisteredTime() time.Time {
	reg, _ := time.Parse(time.RFC3339, user.Registered)
	return reg
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

// GetUserByEmail ...
func GetUserByEmail(email string) (*User, error) {
	rec, err := Get("EmailToUser", email)

	if err != nil {
		return nil, err
	}

	return GetUser(rec["userId"].(string))
}

// AllUsers ...
func AllUsers() (chan *User, error) {
	channel := make(chan *User)
	err := Scan("User", channel)

	return channel, err
}
