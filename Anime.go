package arn

// Anime ...
type Anime struct {
	ID            int                   `json:"id"`
	Type          string                `json:"type"`
	Title         AnimeTitle            `json:"title"`
	Image         string                `json:"image"`
	AiringStatus  string                `json:"airingStatus"`
	Adult         bool                  `json:"adult"`
	Description   string                `json:"description"`
	StartDate     string                `json:"startDate"`
	EndDate       string                `json:"endDate"`
	Hashtag       string                `json:"hashtag"`
	YoutubeID     string                `json:"youtubeId"`
	Source        string                `json:"source"`
	TotalEpisodes int                   `json:"totalEpisodes"`
	Duration      int                   `json:"duration"`
	Watching      int                   `json:"watching"`
	PageGenerated string                `json:"pageGenerated"`
	AnilistEdited uint64                `json:"anilistEdited"`
	Genres        []string              `json:"genres"`
	Tracks        map[string]AnimeTrack `json:"tracks"`
	Links         []AnimeLink           `json:"links"`
	Studios       []AnimeStudio         `json:"studios"`
	Relations     []AnimeRelation       `json:"relations"`

	// "type": "TV",
	// "title": {
	// "romaji": "RWBY",
	// "japanese": "RWBY",
	// "synonyms": [],
	// "english": "RWBY"
	// },
	// "image": "https://static.hummingbird.me/anime/poster_images/000/007/929/large/2e97c8a8938535beb4c86d2e4a50e0cb1376689737_full.jpg",
	// "airingStatus": "finished airing",
	// "adult": 0,
	// "description": "The story takes place in a world filled with monsters and supernatural forces. The series focuses on four girls, each with their own unique weapon and powers, who come together as a team at Beacon Academy in a place called Vale where they are trained to battle these forces alongside other similar teams. Prior to the events of the series, mankind waged a battle of survival against the shadowy creatures of the Grimm until they discovered the power of Dust, which allowed them to fight back the monsters. Dust is used to power magical abilities in the series.",
	// "startDate": "2013-07-18T00:00:00+09:00",
	// "endDate": "2013-11-07T00:00:00+09:00",
	// "hashtag": "#rwby",
	// "youtubeId": "pYW2GmHB5xs",
	// "genres": [
	// "Action",
	// "Adventure",
	// "Fantasy"
	// ],
	// "source": "",
	// "links": [
	// {
	// "url": "http://roosterteeth.com/show/rwby",
	// "title": "Official Site"
	// }
	// ],
	// "studios": [],
	// "relations": [],
	// "totalEpisodes": 16,
	// "duration": 9,
	// "id": 1000001,
	// "created": "2016-08-18T16:37:06.278Z",
	// "createdBy": "4J6qpK1ve",
	// "watching": 0,
	// "pageGenerated": "2016-11-01T14:42:23.461Z",
	// "tracks": {
	// "opening": {
	// "similarity": 0.9411764705882353,
	// "permalink": "https://soundcloud.com/uswatun-hasanah13/rwby-opening-full",
	// "uri": "https://api.soundcloud.com/tracks/189229528",
	// "title": "RWBY Opening full",
	// "likes": 444,
	// "plays": 25502
	// }
	// }
}

// AnimeTitle ...
type AnimeTitle struct {
	Romaji   string   `json:"romaji"`
	English  string   `json:"english"`
	Japanese string   `json:"japanese"`
	Synonyms []string `json:"synonyms"`
}

// AnimeLink ...
type AnimeLink struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// AnimeTrack ...
type AnimeTrack struct {
	URI        string  `json:"uri"`
	Similarity float32 `json:"similarity"`
	Permalink  string  `json:"permalink"`
	Title      string  `json:"title"`
	Likes      int     `json:"likes"`
	Plays      int     `json:"plays"`
}

// AnimeStudio ...
type AnimeStudio struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	IsMainStudio bool   `json:"isMainStudio"`
}

// AnimeRelation ...
type AnimeRelation struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}
