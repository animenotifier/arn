package arn

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ Publishable            = (*Company)(nil)
	_ fmt.Stringer           = (*Company)(nil)
	_ api.Newable            = (*Company)(nil)
	_ api.Editable           = (*Company)(nil)
	_ api.Deletable          = (*Company)(nil)
	_ api.ArrayEventListener = (*Company)(nil)
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
	company.Location = &Location{}

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "create", "Company", company.ID, "", "", "")
	logEntry.Save()

	return company.Unpublish()
}

// Edit saves a log entry for the edit.
func (company *Company) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "edit", "Company", company.ID, key, fmt.Sprint(value.Interface()), fmt.Sprint(newValue.Interface()))
	logEntry.Save()

	return false, nil
}

// OnAppend saves a log entry.
func (company *Company) OnAppend(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayAppend", "Company", company.ID, fmt.Sprintf("%s[%d]", key, index), "", fmt.Sprint(obj))
	logEntry.Save()
}

// OnRemove saves a log entry.
func (company *Company) OnRemove(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayRemove", "Company", company.ID, fmt.Sprintf("%s[%d]", key, index), fmt.Sprint(obj), "")
	logEntry.Save()
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

// DeleteInContext deletes the company in the given context.
func (company *Company) DeleteInContext(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "delete", "Company", company.ID, "", fmt.Sprint(company), "")
	logEntry.Save()

	return company.Delete()
}

// Delete deletes the object from the database.
func (company *Company) Delete() error {
	if company.IsDraft {
		draftIndex := company.Creator().DraftIndex()
		draftIndex.CompanyID = ""
		draftIndex.Save()
	}

	// Remove company ID from all anime
	for anime := range StreamAnime() {
		for index, id := range anime.StudioIDs {
			if id == company.ID {
				anime.StudioIDs = append(anime.StudioIDs[:index], anime.StudioIDs[index+1:]...)
				break
			}
		}

		for index, id := range anime.ProducerIDs {
			if id == company.ID {
				anime.ProducerIDs = append(anime.ProducerIDs[:index], anime.ProducerIDs[index+1:]...)
				break
			}
		}

		for index, id := range anime.LicensorIDs {
			if id == company.ID {
				anime.LicensorIDs = append(anime.LicensorIDs[:index], anime.LicensorIDs[index+1:]...)
				break
			}
		}
	}

	DB.Delete("Company", company.ID)
	return nil
}

// Authorize returns an error if the given API request is not authorized.
func (company *Company) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	if user.Role != "editor" && user.Role != "admin" {
		return errors.New("Insufficient permissions")
	}

	return nil
}
