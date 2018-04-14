package arn

import (
	"errors"
	"os"
	"path"
)

// AMV is an anime music video.
type AMV struct {
	ID            string   `json:"id"`
	File          string   `json:"file" editable:"true" type:"upload" filetype:"video" endpoint:"/api/upload/amv"`
	Title         AMVTitle `json:"title" editable:"true"`
	MainAnimeID   string   `json:"mainAnimeId" editable:"true"`
	ExtraAnimeIDs []string `json:"extraAnimeIds" editable:"true"`
	Tags          []string `json:"tags" editable:"true"`
	IsDraft       bool     `json:"isDraft" editable:"true"`

	HasCreator
	HasEditor
	HasLikes
}

// Link returns the permalink for the AMV.
func (amv *AMV) Link() string {
	return "/amv/" + amv.ID
}

// Publish ...
func (amv *AMV) Publish() error {
	// No draft
	if !amv.IsDraft {
		return errors.New("Not a draft")
	}

	// No anime found
	if amv.MainAnimeID == "" && len(amv.ExtraAnimeIDs) == 0 {
		return errors.New("Need to specify at least one anime")
	}

	draftIndex, err := GetDraftIndex(amv.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.AMVID == "" {
		return errors.New("AMV draft doesn't exist in the user draft index")
	}

	if _, err := os.Stat(path.Join(Root, "videos", "amvs", amv.ID+".webm")); os.IsNotExist(err) {
		return errors.New("You need to upload a WebM file for this AMV")
	}

	amv.IsDraft = false
	draftIndex.AMVID = ""
	draftIndex.Save()

	return nil
}

// Unpublish ...
func (amv *AMV) Unpublish() error {
	draftIndex, err := GetDraftIndex(amv.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.AMVID != "" {
		return errors.New("You still have an unfinished draft")
	}

	amv.IsDraft = true
	draftIndex.AMVID = amv.ID
	draftIndex.Save()
	return nil
}

// OnLike is called when the AMV receives a like.
func (amv *AMV) OnLike(likedBy *User) {
	if likedBy.ID == amv.CreatedBy {
		return
	}

	go func() {
		amv.Creator().SendNotification(&PushNotification{
			Title:   likedBy.Nick + " liked your AMV " + amv.Title.ByUser(amv.Creator()),
			Message: likedBy.Nick + " liked your AMV " + amv.Title.ByUser(amv.Creator()) + ".",
			Icon:    "https:" + likedBy.AvatarLink("large"),
			Link:    "https://notify.moe" + likedBy.Link(),
			Type:    NotificationTypeLike,
		})
	}()
}

// String implements the default string serialization.
func (amv *AMV) String() string {
	return amv.Title.ByUser(nil)
}

// GetAMV returns the AMV with the given ID.
func GetAMV(id string) (*AMV, error) {
	obj, err := DB.Get("AMV", id)

	if err != nil {
		return nil, err
	}

	return obj.(*AMV), nil
}
