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
	c.candidates = c.candidates[1:]
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

func (s stack) peek() (*cellcontext, bool) {
	count := len(s.data)
	if count == 0 {
		return nil, false
	}
	return &s.data[count-1], true
}

func (s *stack) pop() (*cellcontext, bool) {
	value, ok := s.peek()
	if !ok {
		return value, ok
	}
	s.data = s.data[:len(s.data)-1]
	return value, true
}
