package meekstv

import "fmt"

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
	return l.entries[len(l.entries)-2].CandidateSnapshot
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
	CandidateSnapshot []Candidate
	Exhausted         float64
}
