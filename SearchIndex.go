package arn

// SearchIndex ...
type SearchIndex struct {
	TextToID map[string]string `json:"textToId"`
}

// NewSearchIndex ...
func NewSearchIndex() *SearchIndex {
	return &SearchIndex{
		TextToID: make(map[string]string),
	}
}

// GetSearchIndex ...
func GetSearchIndex(id string) (*SearchIndex, error) {
	obj, err := DB.Get("SearchIndex", id)
	return obj.(*SearchIndex), err
}
