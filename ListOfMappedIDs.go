package arn

import (
	"errors"

	"github.com/fatih/color"
)

// ListOfMappedIDs ...
type ListOfMappedIDs struct {
	IDList []*MappedID `json:"idList"`
}

// MappedID ...
type MappedID struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Append appends the given mapped ID to the end of the list.
func (idList *ListOfMappedIDs) Append(typeName string, id string) {
	idList.IDList = append(idList.IDList, &MappedID{
		Type: typeName,
		ID:   id,
	})
}

// Resolve ...
func (idList *ListOfMappedIDs) Resolve() []interface{} {
	var data []interface{}

	for _, mapped := range idList.IDList {
		obj, err := DB.Get(mapped.Type, mapped.ID)

		if err != nil {
			color.Red(err.Error())
			continue
		}

		data = append(data, obj)
	}

	return data
}

// GetListOfMappedIDs ...
func GetListOfMappedIDs(table string, id string) (*ListOfMappedIDs, error) {
	// cache := &ListOfMappedIDs{}
	// err := DB.GetObject(table, id, cache)
	// return cache, err
	return nil, errors.New("Not implemented")
}

// GetForumActivityCached ...
func GetForumActivityCached() ([]Postable, error) {
	cache, err := GetListOfMappedIDs("Cache", "forum activity")

	if err != nil {
		return nil, err
	}

	return ToPostables(cache.Resolve()), nil
}
