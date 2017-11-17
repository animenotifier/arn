package arn

import "github.com/aerogo/nano"

// Analytics ...
type Analytics struct {
	UserID     string              `json:"userId"`
	General    GeneralAnalytics    `json:"general"`
	Screen     ScreenAnalytics     `json:"screen"`
	System     SystemAnalytics     `json:"system"`
	Connection ConnectionAnalytics `json:"connection"`
}

// GeneralAnalytics ...
type GeneralAnalytics struct {
	TimezoneOffset int `json:"timezoneOffset"`
}

// ScreenAnalytics ...
type ScreenAnalytics struct {
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	AvailableWidth  int     `json:"availableWidth"`
	AvailableHeight int     `json:"availableHeight"`
	PixelRatio      float64 `json:"pixelRatio"`
}

// SystemAnalytics ...
type SystemAnalytics struct {
	CPUCount int    `json:"cpuCount"`
	Platform string `json:"platform"`
}

// ConnectionAnalytics ...
type ConnectionAnalytics struct {
	DownLink      float64 `json:"downLink"`
	RoundTripTime float64 `json:"roundTripTime"`
	EffectiveType string  `json:"effectiveType"`
}

// StreamAnalytics returns a stream of all analytics.
func StreamAnalytics() chan *Analytics {
	channel := make(chan *Analytics, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("Analytics") {
			channel <- obj.(*Analytics)
		}

		close(channel)
	}()

	return channel
}

// AllAnalytics returns a slice of all analytics.
func AllAnalytics() []*Analytics {
	var all []*Analytics

	stream := StreamAnalytics()

	for obj := range stream {
		all = append(all, obj)
	}

	return all
}
