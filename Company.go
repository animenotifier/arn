package arn

import (
	"errors"

	"github.com/aerogo/nano"
)

// Company ...
type Company struct {
	ID          string      `json:"id"`
	Name        CompanyName `json:"name" editable:"true"`
	Image       string      `json:"image"`
	Description string      `json:"description" editable:"true" type:"textarea"`
	Mappings    []*Mapping  `json:"mappings" editable:"true"`
	Links       []*Link     `json:"links" editable:"true"`
	Location    Location    `json:"location" editable:"true"`
	Tags        []string    `json:"tags" editable:"true"`
	Likes       []string    `json:"likes"`
	IsDraft     bool        `json:"isDraft"`
	Created     string      `json:"created"`
	CreatedBy   string      `json:"createdBy"`
	Edited      string      `json:"edited"`
	EditedBy    string      `json:"editedBy"`
}

// Link returns a single company.
func (company *Company) Link() string {
	return "/company/" + company.ID
}

// Creator returns the user who created this company.
func (company *Company) Creator() *User {
	user, _ := GetUser(company.CreatedBy)
	return user
}

// Publish ...
func (company *Company) Publish() error {
	// No draft
	if !company.IsDraft {
		return errors.New("Not a draft")
	}

	// No title
	if company.Name.English == "" {
		return errors.New("No English company name")
	}

	draftIndex, err := GetDraftIndex(company.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.CompanyID == "" {
		return errors.New("Company draft doesn't exist in the user draft index")
	}

	company.IsDraft = false
	draftIndex.CompanyID = ""
	draftIndex.Save()
	return nil
}

// Unpublish ...
func (company *Company) Unpublish() error {
	draftIndex, err := GetDraftIndex(company.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.CompanyID != "" {
		return errors.New("You still have an unfinished draft")
	}

	company.IsDraft = true
	draftIndex.CompanyID = company.ID
	draftIndex.Save()
	return nil
}

// Like adds an user to the company's Likes array if they aren't already in it.
func (company *Company) Like(userID string) {
	for _, id := range company.Likes {
		if id == userID {
			return
		}
	}

	company.Likes = append(company.Likes, userID)
}

// Unlike removes the user from the company's Likes array if they are in it.
func (company *Company) Unlike(userID string) {
	for index, id := range company.Likes {
		if id == userID {
			company.Likes = append(company.Likes[:index], company.Likes[index+1:]...)
			return
		}
	}
}

// LikedBy checks to see if the user has liked the company.
func (company *Company) LikedBy(userID string) bool {
	for _, id := range company.Likes {
		if id == userID {
			return true
		}
	}

	return false
}

// GetCompany returns a single company.
func GetCompany(id string) (*Company, error) {
	obj, err := DB.Get("Company", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Company), nil
}

// StreamCompanies returns a stream of all companies.
func StreamCompanies() chan *Company {
	channel := make(chan *Company, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("Company") {
			channel <- obj.(*Company)
		}

		close(channel)
	}()

	return channel
}

// FilterCompanies filters all companies by a custom function.
func FilterCompanies(filter func(*Company) bool) []*Company {
	var filtered []*Company

	channel := DB.All("Company")

	for obj := range channel {
		realObject := obj.(*Company)

		if filter(realObject) {
			filtered = append(filtered, realObject)
		}
	}

	return filtered
}

// AllCompanies returns a slice of all companies.
func AllCompanies() []*Company {
	var all []*Company

	stream := StreamCompanies()

	for obj := range stream {
		all = append(all, obj)
	}

	return all
}
