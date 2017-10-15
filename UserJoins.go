package arn

// Threads ...
func (user *User) Threads() []*Thread {
	threads, _ := GetThreadsByUser(user)
	return threads
}

// Posts ...
func (user *User) Posts() []*Post {
	posts, _ := GetPostsByUser(user)
	return posts
}

// Settings ...
func (user *User) Settings() *Settings {
	if user.settings == nil {
		user.settings, _ = GetSettings(user.ID)
	}

	return user.settings
}

// AnimeList ...
func (user *User) AnimeList() *AnimeList {
	if user.animeList == nil {
		user.animeList, _ = GetAnimeList(user.ID)
	}

	return user.animeList
}

// Follows ...
func (user *User) Follows() *UserFollows {
	if user.follows == nil {
		user.follows, _ = GetUserFollows(user.ID)
	}

	return user.follows
}

// Followers ...
func (user *User) Followers() []*User {
	var followerIDs []string

	for list := range MustStreamUserFollows() {
		if list.Contains(user.ID) {
			followerIDs = append(followerIDs, list.UserID)
		}
	}

	objects, err := DB.GetMany("User", followerIDs)

	if err != nil {
		return nil
	}

	return objects.([]*User)
}

// DraftIndex ...
func (user *User) DraftIndex() *DraftIndex {
	if user.draftIndex == nil {
		user.draftIndex, _ = GetDraftIndex(user.ID)
	}

	return user.draftIndex
}

// SoundTracks returns the soundtracks posted by the user.
func (user *User) SoundTracks() []*SoundTrack {
	tracks, _ := FilterSoundTracks(func(track *SoundTrack) bool {
		return !track.IsDraft && len(track.Media) > 0 && track.CreatedBy == user.ID
	})
	return tracks
}
