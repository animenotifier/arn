package arn

// IgnoreAnimeDifference saves which differences between anime providers can be ignored.
type IgnoreAnimeDifference struct {
	// The ID is built like this: arn:323|mal:356|JapaneseTitle
	ID        string `json:"id"`
	Note      string `json:"note"`
	ValueHash uint64 `json:"valueHash"`
	Created   string `json:"created"`
	CreatedBy string `json:"createdBy"`
}
