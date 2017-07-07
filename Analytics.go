package arn

// Analytics ...
type Analytics struct {
	UserID  string           `json:"userId"`
	General GeneralAnalytics `json:"general"`
	Screen  ScreenAnalytics  `json:"screen"`
	System  SystemAnalytics  `json:"system"`
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

// StreamAnalytics returns a stream of all analytics.
func StreamAnalytics() (chan *Analytics, error) {
	objects, err := DB.All("Analytics")
	return objects.(chan *Analytics), err
}

// MustStreamAnalytics returns a stream of all analytics.
func MustStreamAnalytics() chan *Analytics {
	stream, err := StreamAnalytics()
	PanicOnError(err)
	return stream
}

// AllAnalytics returns a slice of all analytics.
func AllAnalytics() ([]*Analytics, error) {
	var all []*Analytics

	stream, err := StreamAnalytics()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}
