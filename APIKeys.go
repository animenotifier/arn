package arn

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/animenotifier/anilist"
	"github.com/animenotifier/osu"
)

// APIKeys are global API keys for several services
var APIKeys APIKeysData

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

	// FCM struct {
	// 	Key      string `json:"serverKey"`
	// 	SenderID string `json:"senderId"`
	// } `json:"fcm"`

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

	PayPal struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
	} `json:"paypal"`

	VAPID struct {
		Subject    string `json:"subject"`
		PublicKey  string `json:"publicKey"`
		PrivateKey string `json:"privateKey"`
	} `json:"vapid"`
}

func init() {
	root := path.Join(os.Getenv("GOPATH"), "src/github.com/animenotifier/notify.moe")
	apiKeysPath := path.Join(root, "security/api-keys.json")

	if _, err := os.Stat(apiKeysPath); os.IsNotExist(err) {
		defaultAPIKeysPath := path.Join(root, "security/default/api-keys.json")
		err := os.Link(defaultAPIKeysPath, apiKeysPath)

		if err != nil {
			panic(err)
		}
	}

	data, _ := ioutil.ReadFile(apiKeysPath)
	err := json.Unmarshal(data, &APIKeys)

	if err != nil {
		panic(err)
	}

	// Set Osu API key
	osu.APIKey = APIKeys.Osu.Secret

	// Set Anilist API keys
	anilist.APIKeyID = APIKeys.AniList.ID
	anilist.APIKeySecret = APIKeys.AniList.Secret
}
