package arn

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/animenotifier/anilist"
	"github.com/animenotifier/osu"
)

// APIKeys are global API keys for several services
var APIKeys APIKeysData

func init() {
	rootPath := ""
	exe, err := os.Executable()

	if err != nil {
		panic(err)
	}

	if strings.Index(exe, "/notify.moe") == -1 {
		exe, err = os.Getwd()

		if err != nil {
			panic(err)
		}
	}

	arnIndex := strings.Index(exe, "/animenotifier")

	if arnIndex == -1 {
		panic(errors.New("Couldn't find notify.moe directory"))
	} else {
		rootPath = path.Join(exe[:arnIndex], "animenotifier")
	}

	apiKeysPath := path.Join(rootPath, "notify.moe", "security", "api-keys.json")

	if _, err = os.Stat(apiKeysPath); os.IsNotExist(err) {
		// If everything else fails, use hard-coded path.
		// This is needed for some benchmarks and tests.
		apiKeysPath = "/home/eduard/workspace/src/github.com/animenotifier/notify.moe/security/api-keys.json"
	}

	data, _ := ioutil.ReadFile(apiKeysPath)
	err = json.Unmarshal(data, &APIKeys)

	if err != nil {
		panic(err)
	}

	// Set Osu API key
	osu.APIKey = APIKeys.Osu.Secret

	// Set Anilist API keys
	anilist.APIKeyID = APIKeys.AniList.ID
	anilist.APIKeySecret = APIKeys.AniList.Secret
}

// APIKeysData ...
type APIKeysData struct {
	Google struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
	} `json:"google"`

	Facebook struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
	} `json:"facebook"`

	Discord struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
		Token  string `json:"token"`
	} `json:"discord"`

	SoundCloud struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
	} `json:"soundcloud"`

	GoogleAPI struct {
		Key string `json:"key"`
	} `json:"googleAPI"`

	IPInfoDB struct {
		ID string `json:"id"`
	} `json:"ipInfoDB"`

	AniList struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
	} `json:"anilist"`

	Osu struct {
		Secret string `json:"secret"`
	} `json:"osu"`
}
