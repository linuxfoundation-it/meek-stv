package meekstv

import (
	"fmt"
	"strings"
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
		fmt.Printf("transferred surplus: %.02f\n", e.TransferredSurplus)
		fmt.Printf("transferred from elimination: %.02f\n", e.TransferredFromElimination)
		fmt.Printf("exhausted: %.02f\n", e.Exhausted)

		// print the candidate votes in a table containing the name, keep factor, votes, and state
		// for every round
		fmt.Println("candidate\tkeep\tvotes")
		for _, c := range e.CandidateSnapshot {
			fmt.Printf("%s\t%.02f\t%.02f\n", c.Name, c.KeepFactor, c.Votes)
		}

		for _, elected := range e.Elected {
			fmt.Printf("elected %s with %.02f votes\n", elected.Name, elected.Votes)
		}
		for _, defeated := range e.Defeated {
			fmt.Println("eliminating", defeated.Name)
		}
		fmt.Println("-------------------------")
	}
}

func (l *Log) PrintString() string {
	var result strings.Builder

	for _, e := range l.entries {
		result.WriteString(fmt.Sprintf("Round %d:\n", e.Round))
		result.WriteString(fmt.Sprintf("Threshold: %.02f (%.02f%%)\n", e.Threshold, e.Threshold/e.TotVotes*100))
		result.WriteString(fmt.Sprintf("Transferred surplus: %.02f\n", e.TransferredSurplus))
		result.WriteString(fmt.Sprintf("Transferred from elimination: %.02f\n", e.TransferredFromElimination))
		result.WriteString(fmt.Sprintf("Exhausted: %.02f\n", e.Exhausted))

		// print the candidate votes in a table containing the name, keep factor, votes, and state
		result.WriteString("candidate\tkeep\tvotes\n")
		for _, c := range e.CandidateSnapshot {
			result.WriteString(fmt.Sprintf("%s\t%.02f\t%.02f\n", c.Name, c.KeepFactor, c.Votes))
		}

		for _, elected := range e.Elected {
			result.WriteString(fmt.Sprintf("Elected: %s with %.02f votes\n", elected.Name, elected.Votes))
		}
		for _, defeated := range e.Defeated {
			result.WriteString(fmt.Sprintf("Eliminated: %s\n", defeated.Name))
		}
		result.WriteString("-------------------------\n")
	}

	return result.String()
}

func (l *Log) add(round int) {
	l.entries = append(l.entries, &LogEntry{Round: round})
}

func (l *Log) last() *LogEntry {
	return l.entries[len(l.entries)-1]
}

type LogEntry struct {
	Round                      int
	Threshold                  float64
	TotVotes                   float64
	CandidateSnapshot          []Candidate
	Elected                    []Candidate
	Defeated                   []Candidate
	Exhausted                  float64
	TransferredSurplus         float64 // votes transferred from surplus
	TransferredFromElimination float64 // votes transferred from eliminated candidates
}

func (entry *LogEntry) VotesOf(i int) float64 {
	return entry.CandidateSnapshot[i].Votes
}
