package arn

// ListOfIDs ...
type ListOfIDs struct {
	IDList []string `json:"idList"`
}

// Append appends the given ID to the end of the list.
func (idList *ListOfIDs) Append(id string) {
	idList.IDList = append(idList.IDList, id)
}

// GetListOfIDs ...
func GetListOfIDs(table string, id string) (*ListOfIDs, error) {
	cache := &ListOfIDs{}
	err := DB.GetObject(table, id, cache)
	return cache, err
}

// GetAiringAnimeCached ...
func GetAiringAnimeCached() ([]*Anime, error) {
	cache, err := GetListOfIDs("Cache", "airing anime")

	if err != nil {
		return nil, err
	}

	list, err := DB.GetMany("Anime", cache.IDList)

	if err != nil {
		return nil, err
	}

	return list.([]*Anime), nil
}

// GetListOfAnimeCached ...
func GetListOfAnimeCached(cacheKey string) ([]*Anime, error) {
	cache, err := GetListOfIDs("Cache", cacheKey)

	if err != nil {
		return nil, err
	}

	list, err := DB.GetMany("Anime", cache.IDList)

	if err != nil {
		return nil, err
	}

	return list.([]*Anime), nil
}

// GetActiveUsersCached ...
func GetActiveUsersCached() ([]*User, error) {
	cache, err := GetListOfIDs("Cache", "active users")

	if err != nil {
		return nil, err
	}

	list, err := DB.GetMany("User", cache.IDList)

	if err != nil {
		return nil, err
	}

	return list.([]*User), nil
}
