package arn

import "fmt"

// IgnoreAnimeDifference saves which differences between anime providers can be ignored.
type IgnoreAnimeDifference struct {
	// The ID is built like this: arn:323|mal:356|JapaneseTitle
	ID        string `json:"id"`
	ValueHash uint64 `json:"valueHash"`
	Created   string `json:"created"`
	CreatedBy string `json:"createdBy"`
}

// GetIgnoreAnimeDifference ...
func GetIgnoreAnimeDifference(id string) (*IgnoreAnimeDifference, error) {
	obj, err := DB.Get("IgnoreAnimeDifference", id)

	if err != nil {
		return nil, err
	}

	return obj.(*IgnoreAnimeDifference), nil
}

// CreateDifferenceID ...
func CreateDifferenceID(animeID string, dataProvider string, malAnimeID string, typeName string) string {
	return fmt.Sprintf("arn:%s|%s:%s|%s", animeID, dataProvider, malAnimeID, typeName)
}

// IsAnimeDifferenceIgnored tells you whether the given difference is being ignored.
func IsAnimeDifferenceIgnored(animeID string, dataProvider string, malAnimeID string, typeName string, hash uint64) bool {
	key := CreateDifferenceID(animeID, dataProvider, malAnimeID, typeName)
	ignore, err := GetIgnoreAnimeDifference(key)

	if err != nil {
		return false
	}

	return ignore.ValueHash == hash
}
