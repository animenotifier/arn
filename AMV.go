package arn

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/aerogo/nano"
)

// AMV is an anime music video.
type AMV struct {
	ID            string   `json:"id"`
	File          string   `json:"file" editable:"true" type:"upload" filetype:"video" endpoint:"/api/upload/amv/:id/file"`
	Title         AMVTitle `json:"title" editable:"true"`
	MainAnimeID   string   `json:"mainAnimeId" editable:"true"`
	ExtraAnimeIDs []string `json:"extraAnimeIds" editable:"true"`
	Links         []Link   `json:"links" editable:"true"`
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

// SetVideoBytes sets the bytes for the video file.
func (amv *AMV) SetVideoBytes(data []byte) error {
	fileName := amv.ID + ".webm"
	filePath := path.Join(Root, "videos", "amvs", fileName)
	err := ioutil.WriteFile(filePath, data, 0644)

	if err != nil {
		return err
	}

	amv.File = fileName
	return nil
}

// MainAnime returns main anime for the AMV, if available.
func (amv *AMV) MainAnime() *Anime {
	mainAnime, _ := GetAnime(amv.MainAnimeID)
	return mainAnime
}

// ExtraAnime returns main anime for the AMV, if available.
func (amv *AMV) ExtraAnime() []*Anime {
	objects := DB.GetMany("Anime", amv.ExtraAnimeIDs)
	animes := []*Anime{}

	for _, obj := range objects {
		if obj == nil {
			continue
		}

		animes = append(animes, obj.(*Anime))
	}

	return animes
}

// Publish ...
func (amv *AMV) Publish() error {
	// No draft
	if !amv.IsDraft {
		return errors.New("Not a draft")
	}

	// No title
	if amv.Title.String() == "" {
		return errors.New("AMV doesn't have a title")
	}

	// No anime found
	if amv.MainAnimeID == "" && len(amv.ExtraAnimeIDs) == 0 {
		return errors.New("Need to specify at least one anime")
	}

	// No file uploaded
	if amv.File == "" {
		return errors.New("You need to upload a WebM file for this AMV")
	}

	if _, err := os.Stat(path.Join(Root, "videos", "amvs", amv.File)); os.IsNotExist(err) {
		return errors.New("You need to upload a WebM file for this AMV")
	}

	draftIndex, err := GetDraftIndex(amv.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.AMVID == "" {
		return errors.New("AMV draft doesn't exist in the user draft index")
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

// StreamAMVs returns a stream of all AMVs.
func StreamAMVs() chan *AMV {
	channel := make(chan *AMV, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("AMV") {
			channel <- obj.(*AMV)
		}

		close(channel)
	}()

	return channel
}

// AllAMVs returns a slice of all AMVs.
func AllAMVs() []*AMV {
	var all []*AMV

	stream := StreamAMVs()

	for obj := range stream {
		all = append(all, obj)
	}

	return all
}

// FilterAMVs filters all AMVs by a custom function.
func FilterAMVs(filter func(*AMV) bool) []*AMV {
	var filtered []*AMV

	channel := DB.All("AMV")

	for obj := range channel {
		realObject := obj.(*AMV)

		if filter(realObject) {
			filtered = append(filtered, realObject)
		}
	}

	return filtered
}
