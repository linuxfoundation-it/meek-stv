package meekstv

import (
	"testing"

	"github.com/linuxfoundation-it/meek-stv/election"
)

func TestCustomCase_EliminationRecipients(t *testing.T) {
	// choices order defines indices 0..5 matching the provided IDs
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

	// 9 ballots, weight 1 each, zero-based indices as provided
	add(1, []int{5, 2, 4, 0, 1, 3})
	add(1, []int{0, 5, 2, 4, 3, 1})
	add(1, []int{1, 3, 5, 0, 2, 4})
	add(1, []int{0, 4, 1, 3, 5, 2})
	add(1, []int{2, 4, 1, 5, 3, 0})
	add(1, []int{1, 2, 5, 4, 0, 3})
	add(1, []int{1, 2, 5, 0, 3, 4})
	add(1, []int{1, 2, 4, 5, 0, 3})
	add(1, []int{2, 4, 3, 1, 5, 0})

	report := Count(params)

	entries := report.entries
	// Find the round where the previous round eliminated idx 5 (7d72...)
	found := false
	for i := 1; i < len(entries); i++ {
		prev := entries[i-1]
		if len(prev.Defeated) > 0 && prev.Defeated[0].Index == 5 {
			cur := entries[i]
			gotTo468 := cur.EliminationReceived[2]
			gotTo1e6 := cur.EliminationReceived[0]
			if gotTo468 != 1.0 || gotTo1e6 != 0.25 {
				t.Fatalf("unexpected elimination recipients: to 468d=%.2f, to 1e6c=%.2f; want 1.00 and 0.25", gotTo468, gotTo1e6)
			}
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("did not find a round after eliminating idx 5 (7d72...)")
	}
}
