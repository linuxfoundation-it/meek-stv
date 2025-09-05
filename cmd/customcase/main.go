package main

import (
	"fmt"

	"github.com/linuxfoundation-it/meek-stv/election"
	"github.com/linuxfoundation-it/meek-stv/meekstv"
)

func nameByIndex(snapshot []meekstv.Candidate, idx int) string {
	for i := range snapshot {
		if snapshot[i].Index == idx {
			return snapshot[i].Name
		}
	}
	return fmt.Sprintf("candidate#%d", idx)
}

func main() {
	choices := []string{
		"1e6ce7e3-e957-4749-8679-8b2a86751da1", // idx 0
		"2faad174-ea38-47e0-aa47-e0b4b3e146bb", // idx 1
		"468db457-cdee-483c-8182-ad260547524b", // idx 2
		"69399f1d-ee87-496f-ac73-4990af6f91b7", // idx 3
		"6a90af0c-e2bd-48b6-882a-445122e46532", // idx 4
		"7d72a470-19f5-4bea-9c47-6de52b18bf8f", // idx 5
	}

	params := &election.Election{
		Title:          "Custom",
		Candidates:     len(choices),
		Seats:          2,
		Withdrawn:      map[int]bool{},
		Ballots:        []election.Ballot{},
		CandidateNames: choices,
	}

	add := func(weight int, prefs []int) {
		params.Ballots = append(params.Ballots, election.Ballot{Weight: weight, Preferences: prefs})
	}

	add(1, []int{5, 2, 4, 0, 1, 3})
	add(1, []int{0, 5, 2, 4, 3, 1})
	add(1, []int{1, 3, 5, 0, 2, 4})
	add(1, []int{0, 4, 1, 3, 5, 2})
	add(1, []int{2, 4, 1, 5, 3, 0})
	add(1, []int{1, 2, 5, 4, 0, 3})
	add(1, []int{1, 2, 5, 0, 3, 4})
	add(1, []int{1, 2, 4, 5, 0, 3})
	add(1, []int{2, 4, 3, 1, 5, 0})

	report := meekstv.Count(params)

	for i := 0; i < report.NumRounds(); i++ {
		e := report.Round(i)
		fmt.Printf("Round %d:\n", e.Round)
		fmt.Printf("Threshold: %.2f (%.2f%%)\n", e.Threshold, e.Threshold/e.TotVotes*100)
		fmt.Printf("Exhausted: %.2f\n", e.Exhausted)
		fmt.Println("candidate\tkeep\tvotes")
		for _, c := range e.CandidateSnapshot {
			fmt.Printf("%s\t%.2f\t%.2f\n", c.Name, c.KeepFactor, c.Votes)
		}
		if len(e.SurplusReceived) > 0 {
			fmt.Println("surplus transfers:")
			for idx, amt := range e.SurplusReceived {
				fmt.Printf("  -> %s: %.2f\n", nameByIndex(e.CandidateSnapshot, idx), amt)
			}
			if e.SurplusExhaustedDelta > 0 {
				fmt.Printf("  -> exhausted: %.2f\n", e.SurplusExhaustedDelta)
			}
		}
		if len(e.EliminationReceived) > 0 {
			fmt.Println("elimination transfers:")
			for idx, amt := range e.EliminationReceived {
				fmt.Printf("  -> %s: %.2f\n", nameByIndex(e.CandidateSnapshot, idx), amt)
			}
			if e.EliminationExhaustedDelta > 0 {
				fmt.Printf("  -> exhausted: %.2f\n", e.EliminationExhaustedDelta)
			}
		}
		for _, elected := range e.Elected {
			fmt.Printf("Elected: %s with %.2f votes\n", elected.Name, elected.Votes)
		}
		for _, defeated := range e.Defeated {
			fmt.Printf("Eliminated: %s\n", defeated.Name)
		}
		fmt.Println("-------------------------")
	}
}
