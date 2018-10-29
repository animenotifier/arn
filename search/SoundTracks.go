package search

import (
	"sort"
	"strings"

	"github.com/animenotifier/arn"
	"github.com/animenotifier/arn/stringutils"
)

// SoundTracks searches all soundtracks.
func SoundTracks(originalTerm string, maxLength int) []*arn.SoundTrack {
	term := strings.ToLower(stringutils.RemoveSpecialCharacters(originalTerm))

	var results []*Result

	for track := range arn.StreamSoundTracks() {
		if track.ID == originalTerm {
			return []*arn.SoundTrack{track}
		}

		if track.IsDraft {
			continue
		}

		text := strings.ToLower(track.Title.Canonical)
		similarity := stringutils.AdvancedStringSimilarity(term, text)

		if similarity >= MinimumStringSimilarity {
			results = append(results, &Result{
				obj:        track,
				similarity: similarity,
			})
			continue
		}

		text = strings.ToLower(track.Title.Native)
		similarity = stringutils.AdvancedStringSimilarity(term, text)

		if similarity >= MinimumStringSimilarity {
			results = append(results, &Result{
				obj:        track,
				similarity: similarity,
			})
			continue
		}
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].similarity > results[j].similarity
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	// Final list
	final := make([]*arn.SoundTrack, len(results))

	for i, result := range results {
		final[i] = result.obj.(*arn.SoundTrack)
	}

	return final
}
