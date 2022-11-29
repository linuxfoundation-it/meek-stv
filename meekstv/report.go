package meekstv

import (
	"fmt"
)

type Log struct {
	entries []*LogEntry
}

func (l *Log) NumRounds() int {
	return len(l.entries)
}

func (l *Log) Round(i int) *LogEntry {
	if len(l.entries) < i {
		panic(fmt.Errorf("count didn't reach round %d", i))
	}
	return l.entries[i]
}

func (l *Log) Results() []Candidate {
	return l.last().CandidateSnapshot
}

func (l *Log) Winners() []int {
	out := make([]int, 0)
	for _, c := range l.Results() {
		if c.State == Elected {
			out = append(out, c.Index)
		}
	}
	return out
}

func (l *Log) Print() {
	for _, e := range l.entries {
		fmt.Println("round", e.Round)
		fmt.Printf("threshold %.02f (%.02f)\n", e.Threshold, e.Threshold/e.TotVotes*100)

		for _, elected := range e.Elected {
			fmt.Printf("elected %s with %.02f votes\n", elected.Name, elected.Votes)
		}
		for _, defeated := range e.Defeated {
			fmt.Println("eliminating", defeated.Name)
		}
		fmt.Println("-------------------------")
	}
}

func (l *Log) add(round int) {
	l.entries = append(l.entries, &LogEntry{Round: round})
}

func (l *Log) last() *LogEntry {
	return l.entries[len(l.entries)-1]
}

type LogEntry struct {
	Round             int
	Threshold         float64
	TotVotes          float64
	CandidateSnapshot []Candidate
	Elected           []Candidate
	Defeated          []Candidate
	Exhausted         float64
}

func (entry *LogEntry) VotesOf(i int) float64 {
	return entry.CandidateSnapshot[i].Votes
}
