package election

type Election struct {
	Title          string
	Candidates     int
	Seats          int
	Withdrawn      map[int]bool
	Ballots        []Ballot
	CandidateNames []string
}

type Ballot struct {
	Weight      int
	Preferences []int // indices
}
