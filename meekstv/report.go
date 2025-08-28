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
		fmt.Printf("exhausted: %.02f\n", e.Exhausted)

		// print the candidate votes in a table containing the name, keep factor and votes
		// for every round
		fmt.Println("candidate\tkeep\tvotes")
		for _, c := range e.CandidateSnapshot {
			fmt.Printf("%s\t%.02f\t%.02f\n", c.Name, c.KeepFactor, c.Votes)
		}

		// Print transfer breakdowns, if any
		if len(e.SurplusReceived) > 0 {
			fmt.Println("surplus transfers:")
			for idx, amt := range e.SurplusReceived {
				name := e.CandidateSnapshot[idx].Name
				fmt.Printf("  -> %s: %.02f\n", name, amt)
			}
			if e.SurplusExhaustedDelta > 0 {
				fmt.Printf("  -> exhausted: %.02f\n", e.SurplusExhaustedDelta)
			}
		}
		if len(e.EliminationReceived) > 0 {
			fmt.Println("elimination transfers:")
			for idx, amt := range e.EliminationReceived {
				name := e.CandidateSnapshot[idx].Name
				fmt.Printf("  -> %s: %.02f\n", name, amt)
			}
			if e.EliminationExhaustedDelta > 0 {
				fmt.Printf("  -> exhausted: %.02f\n", e.EliminationExhaustedDelta)
			}
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
		result.WriteString(fmt.Sprintf("Exhausted: %.02f\n", e.Exhausted))

		// print the candidate votes in a table containing the name, keep factor and votes
		result.WriteString("candidate\tkeep\tvotes\n")
		for _, c := range e.CandidateSnapshot {
			result.WriteString(fmt.Sprintf("%s\t%.02f\t%.02f\n", c.Name, c.KeepFactor, c.Votes))
		}

		// transfer breakdowns
		if len(e.SurplusReceived) > 0 {
			result.WriteString("surplus transfers:\n")
			for idx, amt := range e.SurplusReceived {
				name := e.CandidateSnapshot[idx].Name
				result.WriteString(fmt.Sprintf("  -> %s: %.02f\n", name, amt))
			}
			if e.SurplusExhaustedDelta > 0 {
				result.WriteString(fmt.Sprintf("  -> exhausted: %.02f\n", e.SurplusExhaustedDelta))
			}
		}
		if len(e.EliminationReceived) > 0 {
			result.WriteString("elimination transfers:\n")
			for idx, amt := range e.EliminationReceived {
				name := e.CandidateSnapshot[idx].Name
				result.WriteString(fmt.Sprintf("  -> %s: %.02f\n", name, amt))
			}
			if e.EliminationExhaustedDelta > 0 {
				result.WriteString(fmt.Sprintf("  -> exhausted: %.02f\n", e.EliminationExhaustedDelta))
			}
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
	Round             int
	Threshold         float64
	TotVotes          float64
	CandidateSnapshot []Candidate
	Elected           []Candidate
	Defeated          []Candidate
	Exhausted         float64

	// Transfer breakdowns realized in this round compared to previous round's event
	// If the previous round elected candidate(s), SurplusReceived shows how much each candidate gained
	// due to surplus redistribution. If the previous round eliminated a candidate, EliminationReceived
	// shows how much each candidate gained from that elimination.
	SurplusReceived           map[int]float64
	EliminationReceived       map[int]float64
	SurplusExhaustedDelta     float64
	EliminationExhaustedDelta float64
}

func (entry *LogEntry) VotesOf(i int) float64 {
	return entry.CandidateSnapshot[i].Votes
}
