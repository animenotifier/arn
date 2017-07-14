package arn

import (
	webpush "github.com/blitzprog/webpush-go"
)

// PushSubscription ...
type PushSubscription struct {
	Platform string `json:"platform"`
	Screen   struct {
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
func (sub *PushSubscription) SendNotification(message string) error {
	s := webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.P256DH,
			Auth:   sub.Auth,
		},
	}

	// Send Notification
	_, err := webpush.SendNotification([]byte(message), &s, &webpush.Options{
		Subscriber:      APIKeys.VAPID.Subject,
		TTL:             60,
		VAPIDPrivateKey: APIKeys.VAPID.PrivateKey,
	})

	return err
}
