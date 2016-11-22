package arn

// ThreadList is a slice of threads. It implements the sort interface.
type ThreadList []*Thread

func (c ThreadList) Len() int {
	return len(c)
}

func (c ThreadList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ThreadList) Less(i, j int) bool {
	a := c[i]
	b := c[j]

	if a.Sticky != b.Sticky {
		if a.Sticky {
			return true
		}

		if b.Sticky {
			return false
		}
	}

	return a.Created > b.Created
}
