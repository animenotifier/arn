package search

import (
	"sort"
	"strings"

	"github.com/animenotifier/arn"
	"github.com/animenotifier/arn/stringutils"
)

// AMVs searches all anime music videos.
func AMVs(originalTerm string, maxLength int) []*arn.AMV {
	term := strings.ToLower(stringutils.RemoveSpecialCharacters(originalTerm))

	var results []*Result

	for amv := range arn.StreamAMVs() {
		if amv.ID == originalTerm {
			return []*arn.AMV{amv}
		}

		if amv.IsDraft {
			continue
		}

		text := strings.ToLower(amv.Title.Canonical)
		similarity := stringutils.AdvancedStringSimilarity(term, text)

		if similarity >= MinimumStringSimilarity {
			results = append(results, &Result{
				obj:        amv,
				similarity: similarity,
			})
			continue
		}

		text = strings.ToLower(amv.Title.Native)
		similarity = stringutils.AdvancedStringSimilarity(term, text)

		if similarity >= MinimumStringSimilarity {
			results = append(results, &Result{
				obj:        amv,
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
	final := make([]*arn.AMV, len(results))

	for i, result := range results {
		final[i] = result.obj.(*arn.AMV)
	}

	return final
}
