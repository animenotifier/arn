package search

import (
	"sort"
	"strings"

	"github.com/animenotifier/arn"
	"github.com/animenotifier/arn/stringutils"
)

// Posts searches all posts.
func Posts(originalTerm string, maxLength int) []*arn.Post {
	term := strings.ToLower(stringutils.RemoveSpecialCharacters(originalTerm))

	var results []*arn.Post

	for post := range arn.StreamPosts() {
		if post.ID == originalTerm {
			return []*arn.Post{post}
		}

		text := strings.ToLower(post.Text)

		if !strings.Contains(text, term) {
			continue
		}

		results = append(results, post)
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].Created > results[j].Created
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	return results
}
