package solver

import "testing"

func TestStack(t *testing.T) {
	c0 := cellcontext{index: 0}
	c1 := cellcontext{index: 1}
	c2 := cellcontext{index: 2}

	s := newStack()
	s.push(c0)
	s.push(c1)
	s.push(c2)

	if !s.hasMore() {
		t.Fatal("stack should not be empty")
	}

	c, ok := s.pop()
	if !ok {
		t.Error("could not pop c2")
	} else if c.index != 2 {
		t.Errorf("after popping c2, expected index 2 but got %d instead", c.index)
	}

	c, ok = s.pop()
	if !ok {
		t.Error("could not pop c1")
	} else if c.index != 1 {
		t.Errorf("after popping c1, expected index 1 but got %d instead", c.index)
	}

	c, ok = s.pop()
	if !ok {
		t.Error("could not pop c0")
	} else if c.index != 0 {
		t.Errorf("after popping c0, expected index 0 but got %d instead", c.index)
	}

	if s.hasMore() {
		t.Error("expected stack to be empty, but it is not")
	}

	c, ok = s.pop()
	if ok {
		t.Error("stack should be empty but it did not return a !ok")
	}
}

func TestCellContext(t *testing.T) {
	tables := []struct {
		hasMore bool
		value   int
	}{
		{true, 0},
		{true, 1},
		{true, 2},
		{false, -1},
	}
	candidates := []int{0, 1, 2}
	c := newCellContext(0, -1, candidates)
	for _, table := range tables {
		if table.hasMore != c.hasMoreCandidates() {
			t.Errorf("expected hasMoreCandidates() to return a value of %v - got %v instead", table.hasMore, c.hasMoreCandidates())
		}
		candidate := c.nextCandidate()
		if candidate != table.value {
			t.Errorf("expected a candidate value of %d - got %d instead", table.value, candidate)
		}
	}
}
