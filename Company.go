package arn

import "github.com/aerogo/nano"

// Company ...
type Company struct {
	ID          string      `json:"id"`
	Name        CompanyName `json:"name"`
	Image       string      `json:"image"`
	Description string      `json:"description" editable:"true"`
	Location    Location    `json:"location"`
	Mappings    []*Mapping  `json:"mappings"`
	IsDraft     bool        `json:"isDraft"`
	Created     string      `json:"created"`
	CreatedBy   string      `json:"createdBy"`
	Edited      string      `json:"edited"`
	EditedBy    string      `json:"editedBy"`
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
