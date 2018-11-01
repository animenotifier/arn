package arn

import (
	"errors"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
	"github.com/animenotifier/arn/autocorrect"
	"github.com/fatih/color"
)

// Force interface implementations
var (
	_ api.Editable = (*User)(nil)
)

// Authorize returns an error if the given API POST request is not authorized.
func (user *User) Authorize(ctx *aero.Context, action string) error {
	editor := GetUserFromContext(ctx)

	if editor == nil {
		return errors.New("Not authorized")
	}

	if editor.ID != ctx.Get("id") && editor.Role != "admin" {
		return errors.New("Can not modify data from other users")
	}

	return nil
}

// Edit updates the user object.
func (user *User) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	// Automatically correct account nicks
	if strings.HasPrefix(key, "Accounts.") && strings.HasSuffix(key, ".Nick") {
		newNick := newValue.String()
		newNick = autocorrect.AccountNick(newNick)
		value.SetString(newNick)

		// Refresh osu info if the name changed
		if key == "Accounts.Osu.Nick" {
			if newNick == "" {
				user.Accounts.Osu.PP = 0
				user.Accounts.Osu.Level = 0
				user.Accounts.Osu.Accuracy = 0
			} else {
				go func() {
					err := user.RefreshOsuInfo()

					if err != nil {
						color.Red("Error refreshing osu info of user '%s' with osu nick '%s': %v", user.Nick, newNick, err)
						return
					}

					color.Green("Refreshed osu info of user '%s' with osu nick '%s': %v", user.Nick, newNick, user.Accounts.Osu.PP)
					user.Save()
				}()
			}
		}

		return true, nil
	}

	// Refresh Overwatch info if the battletag changed
	if key == "Accounts.Overwatch.BattleTag" {
		newBattleTag := newValue.String()
		value.SetString(newBattleTag)

		if newBattleTag == "" {
			user.Accounts.Overwatch.SkillRating = 0
			user.Accounts.Overwatch.Tier = ""
		} else {
			go func() {
				err := user.RefreshOverwatchInfo()

				if err != nil {
					color.Red("Error refreshing Overwatch info of user '%s' with Overwatch battle tag '%s': %v", user.Nick, newBattleTag, err)
					return
				}

				color.Green("Refreshed Overwatch info of user '%s' with Overwatch battle tag '%s': %v", user.Nick, newBattleTag, user.Accounts.Overwatch.SkillRating)
				user.Save()
			}()
		}

		return true, nil
	}

	// Refresh FinalFantasyXIV info if the name or server changed
	if key == "Accounts.FinalFantasyXIV.Nick" || key == "Accounts.FinalFantasyXIV.Server" {
		newValue := newValue.String()
		value.SetString(newValue)

		if newValue == "" {
			user.Accounts.FinalFantasyXIV.Class = ""
			user.Accounts.FinalFantasyXIV.Level = 0
			user.Accounts.FinalFantasyXIV.ItemLevel = 0
		} else if user.Accounts.FinalFantasyXIV.Nick != "" && user.Accounts.FinalFantasyXIV.Server != "" {
			go func() {
				err := user.RefreshFFXIVInfo()

				if err != nil {
					color.Red("Error refreshing FinalFantasy XIV info of user '%s' with nick '%s' on server '%s': %v", user.Nick, user.Accounts.FinalFantasyXIV.Nick, user.Accounts.FinalFantasyXIV.Server, err)
					return
				}

				user.Save()
			}()
		}

		return true, nil
	}

	switch key {
	case "Nick":
		newNick := newValue.String()
		err := user.SetNick(newNick)
		return true, err

	case "Email":
		newEmail := newValue.String()
		err := user.SetEmail(newEmail)
		return true, err

	case "Website":
		newSite := newValue.String()

		if autocorrect.IsTrackerLink(newSite) {
			return true, errors.New("Not an actual website or homepage")
		}

		user.Website = autocorrect.Website(newSite)
		return true, nil

	case "ProExpires":
		user := GetUserFromContext(ctx)

		if user == nil || user.Role != "admin" {
			return true, errors.New("Not authorized to edit")
		}
	}

	return false, nil
}

// Save saves the user object in the database.
func (user *User) Save() {
	DB.Set("User", user.ID, user)
}

// Filter removes privacy critical fields from the user object.
func (user *User) Filter() {
	user.Email = ""
	user.Gender = ""
	user.FirstName = ""
	user.LastName = ""
	user.IP = ""
	user.UserAgent = ""
	user.LastLogin = ""
	user.LastSeen = ""
	user.Accounts.Facebook.ID = ""
	user.Accounts.Google.ID = ""
	user.Accounts.Twitter.ID = ""
	user.AgeRange = UserAgeRange{}
	user.Location = &Location{}
	user.Browser = UserBrowser{}
	user.OS = UserOS{}
}

// ShouldFilter tells whether data needs to be filtered in the given context.
func (user *User) ShouldFilter(ctx *aero.Context) bool {
	ctxUser := GetUserFromContext(ctx)

	if ctxUser != nil && ctxUser.Role == "admin" {
		return false
	}

	return true
}
