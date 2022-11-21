package meekstv

type CandidateState int

const (
	Hopeful CandidateState = iota
	Withdrawn
	Defeated
	Elected
)

type Candidate struct {
	Index      int
	Name       string
	State      CandidateState
	KeepFactor float64
	Votes      float64
	Surplus    float64
}

type Candidates []*Candidate

func (cs Candidates) resetVotes() {
	for _, c := range cs {
		c.Votes = 0.0
	}
}

func (cs Candidates) countState(state CandidateState) int {
	n := 0
	for _, c := range cs {
		if c.State == state {
			n++
		}
	}
	return n
}

func (cs Candidates) countVotes() float64 {
	n := 0.0
	for _, c := range cs {
		n += c.Votes
	}
	return n
}
