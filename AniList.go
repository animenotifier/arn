package arn

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

type anilistListProvider struct {
	AccessToken string
}

// AniList anime provider (singleton)
var AniList = new(anilistListProvider)

// AniListAuthorizeResponse ...
type AniListAuthorizeResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int    `json:"expires"`
	ExpiresIn   int    `json:"expires_in"`
}

// Authorize ...
func (anilist *anilistListProvider) Authorize() error {
	request := gorequest.New().Post("https://anilist.co/api/auth/access_token")

	request.QueryData.Add("grant_type", "client_credentials")
	request.QueryData.Add("client_id", APIKeys.AniList.ID)
	request.QueryData.Add("client_secret", APIKeys.AniList.Secret)

	authorization := &AniListAuthorizeResponse{}
	_, _, errs := request.EndStruct(authorization)

	if len(errs) > 0 {
		return errs[0]
	}

	anilist.AccessToken = authorization.AccessToken

	if anilist.AccessToken == "" {
		return errors.New("Access token is empty")
	}

	return nil
}

// GetAnimeList ...
func (anilist *anilistListProvider) GetAnimeList(user *User) (*AniListAnimeList, error) {
	request := gorequest.New().Get("https://anilist.co/api/user/" + user.Accounts.AniList.Nick + "/animelist?access_token=" + anilist.AccessToken)

	anilistAnimeList := &AniListAnimeList{}
	resp, _, errs := request.EndStruct(anilistAnimeList)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code: %d", resp.StatusCode)
	}

	if len(errs) > 0 {
		return nil, errs[0]
	}

	return anilistAnimeList, nil
}
