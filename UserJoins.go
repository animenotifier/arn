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
	obj, _ := DB.Get("Settings", user.ID)
	return obj.(*Settings)
}

// AnimeList ...
func (user *User) AnimeList() *AnimeList {
	animeList, _ := GetAnimeList(user)
	return animeList
}

// SoundTracks returns the soundtracks posted by the user.
func (user *User) SoundTracks() []*SoundTrack {
	tracks, _ := GetSoundTracksByUser(user)
	return tracks
}
