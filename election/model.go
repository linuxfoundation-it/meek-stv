package election

type Election struct {
	Title          string
	Candidates     int
	Seats          int
	Withdrawn      map[int]bool
	Ballots        []Ballot
	CandidateNames []string
}

func (e *Election) CountEmpty() int {
	n := 0
	for _, b := range e.Ballots {
		if b.IsEmpty() || b.AllWithdrawn(e.Withdrawn) {
			n += b.Weight
		}
	}
	return n
}

type Ballot struct {
	Weight      int
	Preferences []int // indices
}

func (b Ballot) IsEmpty() bool {
	return len(b.Preferences) == 0
}

func (b Ballot) AllWithdrawn(withdrawn map[int]bool) bool {
	a := false
	for _, candidate := range b.Preferences {
		if !withdrawn[candidate] {
			return false
		}
		a = true
	}
	return a
}
