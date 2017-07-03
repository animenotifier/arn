package arn

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
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

// StreamAnime
func (anilist *anilistListProvider) StreamAnime() chan *AniListAnime {
	channel := make(chan *AniListAnime)
	page := 1
	ticker := time.NewTicker(1100 * time.Millisecond)
	rateLimit := ticker.C

	go func() {
		defer close(channel)
		defer ticker.Stop()

		for {
			animePage := []*AniListAnime{}
			request := gorequest.New().Get("https://anilist.co/api/browse/anime?page=" + strconv.Itoa(page) + "&access_token=" + anilist.AccessToken)
			_, _, errs := request.EndStruct(&animePage)

			if len(errs) > 0 {
				color.Red(errs[0].Error())
				page++
				<-rateLimit
				continue
			}

			// We have reached the end
			if len(animePage) == 0 {
				break
			}

			for _, anime := range animePage {
				channel <- anime
			}

			page++
			<-rateLimit
		}
	}()

	return channel
}

// AniListMatch ...
type AniListMatch struct {
	AniListAnime *AniListAnime
	ARNAnime     *Anime
}

// FindAniListAnime tries to find an AniListAnime in our Anime database.
func FindAniListAnime(search *AniListAnime, allAnime []*Anime) *Anime {
	match, err := GetAniListToAnime(strconv.Itoa(search.ID))

	if err == nil {
		anime, _ := GetAnime(match.AnimeID)
		return anime
	}

	if err != nil && !strings.Contains(err.Error(), "not found") {
		color.Red(err.Error())
		return nil
	}

	var mostSimilar *Anime
	var similarity float64

	for _, anime := range allAnime {
		anime.Title.Japanese = strings.Replace(anime.Title.Japanese, "2ndシーズン", "2", 1)
		anime.Title.Romaji = strings.Replace(anime.Title.Romaji, " 2nd Season", " 2", 1)
		search.TitleJapanese = strings.TrimSpace(strings.Replace(search.TitleJapanese, "2ndシーズン", "2", 1))
		search.TitleRomaji = strings.TrimSpace(strings.Replace(search.TitleRomaji, " 2nd Season", " 2", 1))

		titleSimilarity := StringSimilarity(anime.Title.Romaji, search.TitleRomaji)

		if strings.ToLower(anime.Title.Japanese) == strings.ToLower(search.TitleJapanese) {
			titleSimilarity += 1.0
		}

		if strings.ToLower(anime.Title.Romaji) == strings.ToLower(search.TitleRomaji) {
			titleSimilarity += 1.0
		}

		if strings.ToLower(anime.Title.English) == strings.ToLower(search.TitleEnglish) {
			titleSimilarity += 1.0
		}

		if titleSimilarity > similarity {
			mostSimilar = anime
			similarity = titleSimilarity
		}
	}

	if mostSimilar.EpisodeCount != search.TotalEpisodes {
		similarity -= 0.02
	}

	if similarity >= 0.92 && mostSimilar.GetMapping("anilist/anime") == "" {
		// fmt.Printf("MATCH:    %s => %s (%.2f)\n", search.TitleRomaji, mostSimilar.Title.Romaji, similarity)
		mostSimilar.AddMapping("anilist/anime", strconv.Itoa(search.ID), "")
		PanicOnError(mostSimilar.Save())
		return mostSimilar
	}

	// color.Red("MISMATCH: %s => %s (%.2f)", search.TitleRomaji, mostSimilar.Title.Romaji, similarity)

	return nil
}
