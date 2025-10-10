package main

import (
	"fmt"
	"sort"

	"github.com/linuxfoundation-it/meek-stv/election"
	"github.com/linuxfoundation-it/meek-stv/meekstv"
)

func main() {
	fmt.Println("== Scenario A: two ballots (one for each choice) ==")
	runScenario([]election.Ballot{
		{Weight: 1, Preferences: []int{0}},
		{Weight: 1, Preferences: []int{1}},
	})

	fmt.Println("\n== Scenario B: single ballot ranking both choices ==")
	runScenario([]election.Ballot{
		{Weight: 1, Preferences: []int{0, 1}},
	})

	fmt.Println("\n== Scenario C: two ballots both selecting the same choice ==")
	runScenario([]election.Ballot{
		{Weight: 1, Preferences: []int{0, 1}},
		{Weight: 1, Preferences: []int{0, 1}},
	})
}

func runScenario(ballots []election.Ballot) {
	choices := []string{
		"choice-A", // idx 0
		"choice-B", // idx 1
	}

	params := &election.Election{
		Title:          "TwoSeats",
		Candidates:     len(choices),
		Seats:          2,
		Withdrawn:      map[int]bool{},
		Ballots:        ballots,
		CandidateNames: choices,
	}

	report := meekstv.Count(params)

	if report.NumRounds() == 0 {
		fmt.Println("no rounds logged")
		return
	}

	for i := 0; i < report.NumRounds(); i++ {
		e := report.Round(i)
		fmt.Printf("Round %d:\n", e.Round)
		fmt.Printf("Threshold: %.2f (%.2f%%)\n", e.Threshold, safePct(e.Threshold, e.TotVotes))
		fmt.Printf("Exhausted: %.2f\n", e.Exhausted)
		fmt.Println("candidate\tkeep\tvotes")
		for _, c := range e.CandidateSnapshot {
			fmt.Printf("%s\t%.2f\t%.2f\n", c.Name, c.KeepFactor, c.Votes)
		}
		for _, elected := range e.Elected {
			fmt.Printf("Elected: %s with %.2f votes\n", elected.Name, elected.Votes)
		}
		for _, defeated := range e.Defeated {
			fmt.Printf("Eliminated: %s\n", defeated.Name)
		}
		fmt.Println("-------------------------")
	}

	// Final winners summary to verify both are elected
	cands := append([]meekstv.Candidate(nil), report.Results()...)
	sort.Slice(cands, func(i, j int) bool { return cands[i].Votes > cands[j].Votes })
	fmt.Println("Winners:")
	for _, c := range cands {
		if c.State == meekstv.Elected {
			fmt.Printf("- %s (%.2f)\n", c.Name, c.Votes)
		}
	}
}

func safePct(x, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (x / total) * 100
}
