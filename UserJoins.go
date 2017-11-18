package arn

// Threads ...
func (user *User) Threads() []*Thread {
	threads := GetThreadsByUser(user)
	return threads
}

// Posts ...
func (user *User) Posts() []*Post {
	posts, _ := GetPostsByUser(user)
	return posts
}

// Settings ...
func (user *User) Settings() *Settings {
	settings, _ := GetSettings(user.ID)
	return settings
}

// AnimeList ...
func (user *User) AnimeList() *AnimeList {
	animeList, _ := GetAnimeList(user.ID)
	return animeList
}

// Follows ...
func (user *User) Follows() *UserFollows {
	follows, _ := GetUserFollows(user.ID)
	return follows
}

// Followers ...
func (user *User) Followers() []*User {
	var followerIDs []string

	for list := range StreamUserFollows() {
		if list.Contains(user.ID) {
			followerIDs = append(followerIDs, list.UserID)
		}
	}

	usersObj := DB.GetMany("User", followerIDs)
	users := make([]*User, len(usersObj), len(usersObj))

	for i, obj := range usersObj {
		users[i] = obj.(*User)
	}

	return users
}

// FollowersCount ...
func (user *User) FollowersCount() int {
	count := 0

	for list := range StreamUserFollows() {
		if list.Contains(user.ID) {
			count++
		}
	}

	return count
}

// DraftIndex ...
func (user *User) DraftIndex() *DraftIndex {
	draftIndex, _ := GetDraftIndex(user.ID)
	return draftIndex
}

// SoundTracks returns the soundtracks posted by the user.
func (user *User) SoundTracks() []*SoundTrack {
	tracks := FilterSoundTracks(func(track *SoundTrack) bool {
		return !track.IsDraft && len(track.Media) > 0 && track.CreatedBy == user.ID
	})
	return tracks
}
