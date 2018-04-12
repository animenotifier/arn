package arn

import (
	"errors"

	"github.com/aerogo/nano"
)

// Company represents an anime studio, producer or licensor.
type Company struct {
	ID          string      `json:"id"`
	Name        CompanyName `json:"name" editable:"true"`
	Description string      `json:"description" editable:"true" type:"textarea"`
	Links       []*Link     `json:"links" editable:"true"`
	IsDraft     bool        `json:"isDraft"`

	// Mixins
	HasMappings
	HasLikes

	// Other editable fields
	Location *Location `json:"location" editable:"true"`
	Tags     []string  `json:"tags" editable:"true"`

	// Editing dates
	HasCreator
	HasEditor
}

// NewCompany creates a new company.
func NewCompany() *Company {
	return &Company{
		ID:    GenerateID("Company"),
		Name:  CompanyName{},
		Links: []*Link{},
		Tags:  []string{},
		HasCreator: HasCreator{
			Created: DateTimeUTC(),
		},
		HasLikes: HasLikes{
			Likes: []string{},
		},
		HasMappings: HasMappings{
			Mappings: []*Mapping{},
		},
	}
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

// Anime returns the anime connected with this company.
func (company *Company) Anime() (studioAnime []*Anime, producedAnime []*Anime, licensedAnime []*Anime) {
	for anime := range StreamAnime() {
		if Contains(anime.StudioIDs, company.ID) {
			studioAnime = append(studioAnime, anime)
		}

		if Contains(anime.ProducerIDs, company.ID) {
			producedAnime = append(producedAnime, anime)
		}

		if Contains(anime.LicensorIDs, company.ID) {
			licensedAnime = append(licensedAnime, anime)
		}
	}

	SortAnimeByQuality(studioAnime)
	SortAnimeByQuality(producedAnime)
	SortAnimeByQuality(licensedAnime)

	return studioAnime, producedAnime, licensedAnime
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

// String implements the default string serialization.
func (company *Company) String() string {
	return company.Name.English
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
