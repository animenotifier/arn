package arn

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	webpush "github.com/blitzprog/webpush-go"
)

// PushSubscription ...
type PushSubscription struct {
	Platform  string `json:"platform"`
	UserAgent string `json:"userAgent"`
	Screen    struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"screen"`
	Endpoint string `json:"endpoint"`
	P256DH   string `json:"p256dh"`
	Auth     string `json:"auth"`
	Created  string `json:"created"`
}

// ID ...
func (sub *PushSubscription) ID() string {
	return sub.Endpoint
}

// SendNotification ...
func (sub *PushSubscription) SendNotification(notification *PushNotification) error {
	// Define endpoint and security tokens
	s := webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.P256DH,
			Auth:   sub.Auth,
		},
	}

	// Create notification
	data, err := json.Marshal(notification)

	if err != nil {
		return err
	}

	// Send Notification
	resp, err := webpush.SendNotification(data, &s, &webpush.Options{
		Subscriber:      APIKeys.VAPID.Subject,
		TTL:             60,
		VAPIDPrivateKey: APIKeys.VAPID.PrivateKey,
	})

	if err != nil {
		return err
	}

	// Return "Subscription expired" so that it's marked for deletion by the caller
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusGone {
			return errors.New("Subscription expired")
		}

		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body))
	}

	return nil
}
