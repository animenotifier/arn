package arn

// AnimeByPopularity is a slice of anime. It implements the sort interface.
type AnimeByPopularity []*Anime

func (c AnimeByPopularity) Len() int {
	return len(c)
}

func (c AnimeByPopularity) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c AnimeByPopularity) Less(i, j int) bool {
	return c[i].Watching > c[j].Watching
}
