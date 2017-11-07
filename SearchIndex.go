package arn

// // SearchIndex ...
// type SearchIndex struct {
// 	TextToID map[string]string `json:"textToId"`
// }

// // NewSearchIndex ...
// func NewSearchIndex() *SearchIndex {
// 	return &SearchIndex{
// 		TextToID: make(map[string]string),
// 	}
// }

// // GetSearchIndex ...
// func GetSearchIndex(id string) (*SearchIndex, error) {
// 	obj, err := DB.Get("SearchIndex", id)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return obj.(*SearchIndex), nil
// }

// // Search2 is a fuzzy search.
// func Search2(term string, maxUsers, maxAnime, maxPosts, maxThreads int) ([]*User, []*Anime, []*Post, []*Thread) {
// 	term = strings.ToLower(term)

// 	if term == "" {
// 		return nil, nil, nil, nil
// 	}

// 	var userResults []*User
// 	var animeResults []*Anime
// 	var postResults []*Post
// 	var threadResults []*Thread

// 	type SearchItem struct {
// 		text       string
// 		similarity float64
// 	}

// 	searchUsers := func() {
// 		// Search userResults
// 		var user *User

// 		userSearchIndex, err := GetSearchIndex("User")

// 		if err != nil {
// 			return
// 		}

// 		textToID := userSearchIndex.TextToID

// 		// Search items
// 		items := make([]*SearchItem, 0)

// 		for name := range textToID {
// 			s := StringSimilarity(term, name)

// 			if strings.Contains(name, term) {
// 				s += 0.5
// 			}

// 			if s < MinimumStringSimilarity {
// 				continue
// 			}

// 			items = append(items, &SearchItem{
// 				text:       name,
// 				similarity: s,
// 			})
// 		}

// 		// Sort
// 		sort.Slice(items, func(i, j int) bool {
// 			return items[i].similarity > items[j].similarity
// 		})

// 		// Limit
// 		if len(items) >= maxUsers {
// 			items = items[:maxUsers]
// 		}

// 		// Fetch data
// 		for _, item := range items {
// 			user, err = GetUser(textToID[item.text])

// 			if err != nil {
// 				continue
// 			}

// 			userResults = append(userResults, user)
// 		}
// 	}

// 	searchAnime := func() {
// 		// Remove special characters when searching anime titles
// 		animeSearchTerm := RemoveSpecialCharacters(term)

// 		// Search anime
// 		var anime *Anime

// 		animeSearchIndex, err := GetSearchIndex("Anime")

// 		if err != nil {
// 			return
// 		}

// 		textToID := animeSearchIndex.TextToID

// 		// Search items
// 		items := make([]*SearchItem, 0)
// 		animeIDAdded := map[string]*SearchItem{}

// 		for name, id := range textToID {
// 			cleanName := RemoveSpecialCharacters(name)
// 			s := StringSimilarity(animeSearchTerm, cleanName)

// 			if strings.Contains(cleanName, animeSearchTerm) {
// 				s += 0.5
// 			}

// 			if s < MinimumStringSimilarity {
// 				continue
// 			}

// 			addedEntry, found := animeIDAdded[id]

// 			// Skip existing anime IDs
// 			if found {
// 				// But update existing entry with new similarity if it's higher
// 				if s > addedEntry.similarity {
// 					addedEntry.similarity = s
// 				}

// 				continue
// 			}

// 			item := &SearchItem{
// 				text:       name,
// 				similarity: s,
// 			}
// 			items = append(items, item)

// 			animeIDAdded[id] = item
// 		}

// 		// Sort
// 		sort.Slice(items, func(i, j int) bool {
// 			return items[i].similarity > items[j].similarity
// 		})

// 		// Limit
// 		if len(items) >= maxAnime {
// 			items = items[:maxAnime]
// 		}

// 		// Fetch data
// 		for _, item := range items {
// 			anime, err = GetAnime(textToID[item.text])

// 			if err != nil {
// 				continue
// 			}

// 			animeResults = append(animeResults, anime)
// 		}
// 	}

// 	searchPosts := func() {
// 		postSearchIndex, err := GetSearchIndex("Post")

// 		if err != nil {
// 			return
// 		}

// 		textToID := postSearchIndex.TextToID

// 		// Search items
// 		items := make([]string, 0)

// 		for text, postID := range textToID {
// 			if !strings.Contains(text, term) {
// 				continue
// 			}

// 			items = append(items, postID)

// 			// Limit
// 			if len(items) >= maxPosts {
// 				break
// 			}
// 		}

// 		// Fetch data
// 		objects := DB.GetMany("Post", items)
// 		postResults = make([]*Post, len(objects), len(objects))

// 		for i, obj := range objects {
// 			postResults[i] = obj.(*Post)
// 		}
// 	}

// 	searchThreads := func() {
// 		threadSearchIndex, err := GetSearchIndex("Thread")

// 		if err != nil {
// 			return
// 		}

// 		textToID := threadSearchIndex.TextToID

// 		// Search items
// 		items := make([]string, 0)

// 		for text, threadID := range textToID {
// 			if !strings.Contains(text, term) {
// 				continue
// 			}

// 			items = append(items, threadID)

// 			// Limit
// 			if len(items) >= maxThreads {
// 				break
// 			}
// 		}

// 		// Fetch data
// 		objects := DB.GetMany("Thread", items)
// 		threadResults = make([]*Thread, len(objects), len(objects))

// 		for i, obj := range objects {
// 			threadResults[i] = obj.(*Thread)
// 		}
// 	}

// 	// Search everything in parallel
// 	flow.Parallel(
// 		searchUsers,
// 		searchAnime,
// 		searchPosts,
// 		searchThreads,
// 	)

// 	return userResults, animeResults, postResults, threadResults
// }
