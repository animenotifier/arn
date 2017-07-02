package arn

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

var setNickMutex sync.Mutex

// User ...
type User struct {
	ID         string       `json:"id"`
	Nick       string       `json:"nick" editable:"true"`
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
	Tagline    string       `json:"tagline" editable:"true"`
	Website    string       `json:"website" editable:"true"`
	IP         string       `json:"ip"`
	UserAgent  string       `json:"agent"`
	AgeRange   UserAgeRange `json:"ageRange"`
	Location   UserLocation `json:"location"`
	Accounts   UserAccounts `json:"accounts"`
	Browser    UserBrowser  `json:"browser"`
	OS         UserOS       `json:"os"`
	Following  []string     `json:"following"`
}

// NewUser creates an empty user object with a unique ID.
func NewUser() *User {
	user := &User{
		ID: GenerateID("User"),

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

// Link returns the URI to the user page.
func (user *User) Link() string {
	return "/+" + user.Nick
}

// CoverImageURL ...
func (user *User) CoverImageURL() string {
	return "/images/cover/default"
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

// SetNick changes the user's nickname safely.
func (user *User) SetNick(newName string) error {
	setNickMutex.Lock()
	defer setNickMutex.Unlock()

	newName = FixUserNick(newName)

	if !IsValidNick(newName) {
		return errors.New("Invalid nickname")
	}

	if newName == user.Nick {
		return nil
	}

	// Make sure the nickname doesn't exist already
	_, err := GetUserByNick(newName)

	// If there was no error: the username exists.
	// If "not found" is not included in the error message it's a different error type.
	if err == nil || strings.Index(err.Error(), "not found") == -1 {
		return errors.New("Username '" + newName + "' is taken already")
	}

	return user.ForceSetNick(newName)
}

// ForceSetNick forces a nickname overwrite.
func (user *User) ForceSetNick(newName string) error {
	// Delete old nick reference
	DB.Delete("NickToUser", user.Nick)

	// Set new nick
	user.Nick = newName

	return DB.Set("NickToUser", user.Nick, &NickToUser{
		Nick:   user.Nick,
		UserID: user.ID,
	})
}

// RefreshOsuInfo refreshes a user's Osu information.
func (user *User) RefreshOsuInfo() error {
	if user.Accounts.Osu.Nick == "" {
		return nil
	}

	osu, err := GetOsuUser(user.Accounts.Osu.Nick)

	if err != nil {
		return err
	}

	user.Accounts.Osu.PP, _ = strconv.ParseFloat(osu.PPRaw, 64)
	user.Accounts.Osu.Level, _ = strconv.ParseFloat(osu.Level, 64)
	user.Accounts.Osu.Accuracy, _ = strconv.ParseFloat(osu.Accuracy, 64)

	return nil
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
