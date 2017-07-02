package arn

import (
	"errors"
	"strconv"

	"github.com/parnurzeal/gorequest"
)

// OsuUser is a user in Osu.
type OsuUser struct {
	UserID        string        `json:"user_id"`
	UserName      string        `json:"username"`
	Count300      string        `json:"count300"`
	Count100      string        `json:"count100"`
	Count50       string        `json:"count50"`
	PlayCount     string        `json:"playcount"`
	RankedScore   string        `json:"ranked_score"`
	TotalScore    string        `json:"total_score"`
	PPRank        string        `json:"pp_rank"`
	Level         string        `json:"level"`
	PPRaw         string        `json:"pp_raw"`
	Accuracy      string        `json:"accuracy"`
	CountryRankSS string        `json:"count_rank_ss"`
	CountryRankS  string        `json:"count_rank_s"`
	CountryRankA  string        `json:"count_rank_a"`
	Country       string        `json:"country"`
	PPCountryRank string        `json:"pp_country_rank"`
	Events        []interface{} `json:"events"`
}

// GetOsuUser ...
func GetOsuUser(nick string) (*OsuUser, error) {
	users := []*OsuUser{}

	request := gorequest.New().Get("https://osu.ppy.sh/api/get_user?u=" + nick + "&type=string&k=" + APIKeys.Osu.Secret)
	request = request.Param("Accept", "application/json")

	_, _, errs := request.EndStruct(&users)

	if len(errs) > 0 {
		return nil, errs[0]
	}

	if len(users) == 0 {
		return nil, errors.New("User not found")
	}

	return users[0], nil
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
