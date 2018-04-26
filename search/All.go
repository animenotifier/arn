package search

import (
	"github.com/aerogo/flow"
	"github.com/animenotifier/arn"
)

// MinimumStringSimilarity is the minimum JaroWinkler distance we accept for search results.
const MinimumStringSimilarity = 0.89

// popularityDamping reduces the factor of popularity in search results.
const popularityDamping = 0.001

// Result ...
type Result struct {
	obj        interface{}
	similarity float64
}

// All is a fuzzy search.
func All(term string, maxUsers, maxAnime, maxPosts, maxThreads, maxTracks, maxCharacters, maxCompanies int) ([]*arn.User, []*arn.Anime, []*arn.Post, []*arn.Thread, []*arn.SoundTrack, []*arn.Character, []*arn.Company) {
	if term == "" {
		return nil, nil, nil, nil, nil, nil, nil
	}

	var userResults []*arn.User
	var animeResults []*arn.Anime
	var postResults []*arn.Post
	var threadResults []*arn.Thread
	var trackResults []*arn.SoundTrack
	var characterResults []*arn.Character
	var companyResults []*arn.Company

	flow.Parallel(func() {
		userResults = Users(term, maxUsers)
	}, func() {
		animeResults = Anime(term, maxAnime)
	}, func() {
		postResults = Posts(term, maxPosts)
	}, func() {
		threadResults = Threads(term, maxThreads)
	}, func() {
		trackResults = SoundTracks(term, maxTracks)
	}, func() {
		characterResults = Characters(term, maxCharacters)
	}, func() {
		companyResults = Companies(term, maxCompanies)
	})

	return userResults, animeResults, postResults, threadResults, trackResults, characterResults, companyResults
}
