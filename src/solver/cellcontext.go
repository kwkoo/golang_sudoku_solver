package solver

type cellcontext struct {
	index      int
	nextEmpty  int
	candidates []int
}

func newCellContext(index, nextEmpty int, candidates []int) cellcontext {
	return cellcontext{
		index:      index,
		nextEmpty:  nextEmpty,
		candidates: candidates,
	}
}

func (c cellcontext) hasMoreCandidates() bool {
	return len(c.candidates) > 0
}

// returns -1 if there are no more candidates
func (c *cellcontext) nextCandidate() int {
	candidateCount := len(c.candidates)
	if candidateCount == 0 {
		return -1
	}
	value := c.candidates[0]
	if candidateCount == 1 {
		c.candidates = []int{}
	} else {
		c.candidates = c.candidates[1:]
	}
	return value
}

type stack struct {
	data []cellcontext
}

func newStack() stack {
	return stack{data: []cellcontext{}}
}

func (s stack) hasMore() bool {
	return len(s.data) > 0
}

func (s *stack) push(c cellcontext) {
	s.data = append(s.data, c)
}

func (s *stack) pop() (cellcontext, bool) {
	count := len(s.data)
	if count == 0 {
		return cellcontext{}, false
	}
	var value cellcontext
	if count == 1 {
		value = s.data[0]
		s.data = []cellcontext{}

	} else {
		value = s.data[count-1]
		s.data = s.data[:count-1]
	}
	return value, true
}
