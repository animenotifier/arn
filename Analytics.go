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
