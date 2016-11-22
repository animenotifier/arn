package arn

// PostList is a slice of posts. It implements the sort interface.
type PostList []*Post

func (c PostList) Len() int {
	return len(c)
}

func (c PostList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c PostList) Less(i, j int) bool {
	return c[i].Created < c[j].Created
}
