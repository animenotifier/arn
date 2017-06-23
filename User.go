package arn

import (
	"errors"
	"strings"
	"time"

	shortid "github.com/ventu-io/go-shortid"
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
	LastSeen   string       `json:"lastSeen"`
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

// UserAgeRange ...
type UserAgeRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
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

// GoogleToUser ...
type GoogleToUser struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

// FacebookToUser is the same data structure as GoogleToUser
type FacebookToUser GoogleToUser

// TwitterToUser is the same data structure as GoogleToUser
type TwitterToUser GoogleToUser

// NewUser creates an empty user object with a unique ID.
func NewUser() *User {
	user := &User{
		ID: GenerateUserID(),

		// Avoid nil value fields
		Following: make([]string, 0),
	}

	return user
}

// RegisterUser registers a new user in the database and sets up all the required references.
func RegisterUser(user *User) error {
	var err error

	user.Registered = DateTimeUTC()
	user.LastLogin = user.Registered
	user.LastSeen = user.Registered

	// Save nick in NickToUser table
	err = DB.Set("NickToUser", user.Nick, &NickToUser{
		Nick:   user.Nick,
		UserID: user.ID,
	})

	if err != nil {
		return err
	}

	// Save email in EmailToUser table
	err = DB.Set("EmailToUser", user.Email, &EmailToUser{
		Email:  user.Email,
		UserID: user.ID,
	})

	if err != nil {
		return err
	}

	// Create default settings
	err = NewSettings(user.ID).Save()

	if err != nil {
		return err
	}

	// Add empty anime list
	err = DB.Set("AnimeList", user.ID, &AnimeList{
		UserID: user.ID,
		Items:  make([]*AnimeListItem, 0),
	})

	if err != nil {
		return err
	}

	return nil
}

// ConnectGoogle connects the user's account with a Google account.
func (user *User) ConnectGoogle(googleID string) error {
	user.Accounts.Google.ID = googleID

	return DB.Set("GoogleToUser", googleID, &GoogleToUser{
		ID:     googleID,
		UserID: user.ID,
	})
}

// ConnectFacebook connects the user's account with a Facebook account.
func (user *User) ConnectFacebook(facebookID string) error {
	user.Accounts.Facebook.ID = facebookID

	return DB.Set("FacebookToUser", facebookID, &FacebookToUser{
		ID:     facebookID,
		UserID: user.ID,
	})
}

// ConnectTwitter connects the user's account with a Twitter account.
func (user *User) ConnectTwitter(twtterID string) error {
	user.Accounts.Twitter.ID = twtterID

	return DB.Set("TwitterToUser", twtterID, &TwitterToUser{
		ID:     twtterID,
		UserID: user.ID,
	})
}

// GenerateUserID generates a unique user ID.
func GenerateUserID() string {
	id, _ := shortid.Generate()

	// Retry until we find an unused ID
	retry := 0

	for {
		_, err := GetUser(id)

		if err != nil && strings.Index(err.Error(), "not found") != -1 {
			return id
		}

		retry++

		if retry > 10 {
			panic(errors.New("Can't generate unique user ID"))
		}

		id, _ = shortid.Generate()
	}
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

// HasAvatar ...
func (user *User) HasAvatar() bool {
	return user.Avatar != ""
}

// SmallAvatar ...
func (user *User) SmallAvatar() string {
	return "//media.notify.moe/images/avatars/small/" + user.ID + ".webp"
}

// LargeAvatar ...
func (user *User) LargeAvatar() string {
	return "//media.notify.moe/images/avatars/large/" + user.ID + ".webp"
}

// RealName returns the real name of the user.
func (user *User) RealName() string {
	if user.LastName == "" {
		return user.FirstName
	}

	if user.FirstName == "" {
		return user.LastName
	}

	return user.FirstName + " " + user.LastName
}

// Settings ...
func (user *User) Settings() *Settings {
	obj, _ := DB.Get("Settings", user.ID)
	return obj.(*Settings)
}

// AnimeList ...
func (user *User) AnimeList() *AnimeList {
	animeList, _ := GetAnimeList(user.ID)
	return animeList
}

// Save saves the user object in the database.
func (user *User) Save() error {
	return DB.Set("User", user.ID, user)
}

// SetNick changes the user's nickname safely.
func (user *User) SetNick(newName string) error {
	if !IsValidNick(user.Nick) {
		return errors.New("Invalid nickname")
	}

	// Delete old nick reference
	DB.Delete("NickToUser", user.Nick)

	// Set new nick
	user.Nick = newName

	return DB.Set("NickToUser", user.Nick, &NickToUser{
		Nick:   user.Nick,
		UserID: user.ID,
	})
}

// SetEmail changes the user's email safely.
func (user *User) SetEmail(newName string) error {
	if !IsValidEmail(user.Email) {
		return errors.New("Invalid email address")
	}

	// Delete old email reference
	DB.Delete("EmailToUser", user.Email)

	// Set new email
	user.Email = newName

	return DB.Set("EmailToUser", user.Email, &EmailToUser{
		Email:  user.Email,
		UserID: user.ID,
	})
}

// RegisteredTime ...
func (user *User) RegisteredTime() time.Time {
	reg, _ := time.Parse(time.RFC3339, user.Registered)
	return reg
}

// IsActive ...
func (user *User) IsActive() bool {
	// Exclude people who didn't change their nickname.
	if strings.HasPrefix(user.Nick, "g") || strings.HasPrefix(user.Nick, "fb") || strings.HasPrefix(user.Nick, "t") {
		return false
	}

	return true
}

// WebsiteURL adds https:// to the URL.
func (user *User) WebsiteURL() string {
	return "https://" + user.Website
}

// Threads ...
func (user *User) Threads() []*Thread {
	threads, _ := GetThreadsByUser(user)
	return threads
}

// Posts ...
func (user *User) Posts() []*Post {
	posts, _ := GetPostsByUser(user)
	return posts
}

// Link returns the URI to the user page.
func (user *User) Link() string {
	return "/+" + user.Nick
}

// GetUser ...
func GetUser(id string) (*User, error) {
	obj, err := DB.Get("User", id)
	return obj.(*User), err
}

// GetUserByNick ...
func GetUserByNick(nick string) (*User, error) {
	return GetUserFromTable("NickToUser", nick)
}

// GetUserByEmail ...
func GetUserByEmail(email string) (*User, error) {
	return GetUserFromTable("EmailToUser", email)
}

// GetUserFromTable ...
func GetUserFromTable(table string, id string) (*User, error) {
	rec, err := DB.GetMap(table, id)

	if err != nil {
		return nil, err
	}

	return GetUser(rec["userId"].(string))
}

// AllUsers ...
func AllUsers() (chan *User, error) {
	channel := make(chan *User)
	err := DB.Scan("User", channel)

	return channel, err
}

// FilterUsers filters all users by a custom function.
func FilterUsers(filter func(*User) bool) ([]*User, error) {
	var filtered []*User

	channel, err := AllUsers()

	if err != nil {
		return filtered, err
	}

	for obj := range channel {
		if filter(obj) {
			filtered = append(filtered, obj)
		}
	}

	return filtered, nil
}
