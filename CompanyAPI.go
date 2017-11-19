package arn

import (
	"errors"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ Publishable  = (*Company)(nil)
	_ api.Newable  = (*Company)(nil)
	_ api.Editable = (*Company)(nil)
)

// Actions
func init() {
	API.RegisterActions("Company", []*api.Action{
		// Publish
		PublishAction(),

		// Unpublish
		UnpublishAction(),

		// Like
		LikeAction(),

		// Unlike
		UnlikeAction(),
	})
}

// Create sets the data for a new company with data we received from the API request.
func (company *Company) Create(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	company.ID = GenerateID("Company")
	company.Created = DateTimeUTC()
	company.CreatedBy = user.ID
	company.Mappings = []*Mapping{}
	company.Links = []*Link{}
	company.Tags = []string{}
	company.Likes = []string{}

	return company.Unpublish()
}

// AfterEdit updates the metadata.
func (company *Company) AfterEdit(ctx *aero.Context) error {
	company.Edited = DateTimeUTC()
	company.EditedBy = GetUserFromContext(ctx).ID
	return nil
}

// Save saves the company in the database.
func (company *Company) Save() {
	DB.Set("Company", company.ID, company)
}

// Delete deletes the object from the database.
func (company *Company) Delete() error {
	if company.IsDraft {
		draftIndex := company.Creator().DraftIndex()
		draftIndex.CompanyID = ""
		draftIndex.Save()
	}

	DB.Delete("Company", company.ID)
	return nil
}

// Authorize returns an error if the given API request is not authorized.
func (company *Company) Authorize(ctx *aero.Context, action string) error {
	return AuthorizeIfLoggedIn(ctx)
}
