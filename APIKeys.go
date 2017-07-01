package arn

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// APIKeys are global API keys for several services
var APIKeys APIKeysData

func init() {
	exe, err := os.Executable()

	if err != nil {
		panic(err)
	}

	notifyMoeIndex := strings.Index(exe, "notify.moe")

	if notifyMoeIndex == -1 {
		panic(errors.New("Couldn't find notify.moe directory"))
	}

	rootPath := exe[:notifyMoeIndex]

	data, _ := ioutil.ReadFile(path.Join(rootPath, "notify.moe", "security", "api-keys.json"))
	err = json.Unmarshal(data, &APIKeys)

	if err != nil {
		panic(err)
	}
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
}
