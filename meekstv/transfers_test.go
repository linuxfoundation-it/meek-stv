package meekstv

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/linuxfoundation-it/meek-stv/election"
)

type expectedJSON struct {
	Winners []int  `json:"winners"`
	Title   string `json:"title"`
}

func readJSONExpect(t *testing.T, jsonPath string) expectedJSON {
	t.Helper()
	f, err := os.Open(jsonPath)
	if err != nil {
		t.Fatalf("open json: %v", err)
	}
	defer f.Close()
	var exp expectedJSON
	if err := json.NewDecoder(f).Decode(&exp); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	return exp
}

func readElectionTxt(t *testing.T, txtPath string) *election.Election {
	t.Helper()
	f, err := os.Open(txtPath)
	if err != nil {
		t.Fatalf("open txt: %v", err)
	}
	defer f.Close()
	return election.Read(f)
}

func floatAlmostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func asSortedCopy(ints []int) []int {
	out := append([]int(nil), ints...)
	sort.Ints(out)
	return out
}

func TestElection10_WinnersAndTransfers(t *testing.T) {
	base := filepath.Join("..", "testdata")
	txt := filepath.Join(base, "election10.txt")
	jsonPath := filepath.Join(base, "election10.json")

	params := readElectionTxt(t, txt)
	report := Count(params)

	exp := readJSONExpect(t, jsonPath)

	// winners match
	gotWinners := report.Winners()
	if !slicesEqual(asSortedCopy(gotWinners), asSortedCopy(exp.Winners)) {
		t.Fatalf("winners mismatch: got %v, want %v", gotWinners, exp.Winners)
	}

	// transfer conservation per round
	checkTransferConservation(t, &report)
}

func TestElection11_WinnersAndTransfers(t *testing.T) {
	base := filepath.Join("..", "testdata")
	txt := filepath.Join(base, "election11.txt")
	jsonPath := filepath.Join(base, "election11.json")

	params := readElectionTxt(t, txt)
	report := Count(params)

	exp := readJSONExpect(t, jsonPath)

	// winners match
	gotWinners := report.Winners()
	if !slicesEqual(asSortedCopy(gotWinners), asSortedCopy(exp.Winners)) {
		t.Fatalf("winners mismatch: got %v, want %v", gotWinners, exp.Winners)
	}

	// transfer conservation per round
	checkTransferConservation(t, &report)
}

func checkTransferConservation(t *testing.T, report *Log) {
	t.Helper()
	entries := report.entries
	if len(entries) < 2 {
		return
	}
	const tol = 1e-2
	for i := 1; i < len(entries); i++ {
		prev := entries[i-1]
		cur := entries[i]
		// If previous round eliminated someone, positive deltas + exhausted delta should equal eliminated candidate's prior votes
		if len(prev.Defeated) > 0 {
			defeatedIdx := prev.Defeated[0].Index
			defeatedVotes := prev.CandidateSnapshot[defeatedIdx].Votes
			sumRecipients := 0.0
			for _, v := range cur.EliminationReceived {
				sumRecipients += v
			}
			total := sumRecipients + cur.EliminationExhaustedDelta
			if !floatAlmostEqual(total, defeatedVotes, tol) {
				t.Fatalf("round %d elimination conservation failed: got %.2f (recipients %.2f + exhausted %.2f), want %.2f", cur.Round, total, sumRecipients, cur.EliminationExhaustedDelta, defeatedVotes)
			}
		}
		// If previous round elected someone, positive deltas + exhausted delta should equal the actual drop in elected candidates' votes between rounds
		if len(prev.Elected) > 0 {
			// compute total drop among candidates who were elected in prev round
			prevSnap := prev.CandidateSnapshot
			curSnap := cur.CandidateSnapshot
			totalDrop := 0.0
			for _, ec := range prev.Elected {
				idx := ec.Index
				if idx >= 0 && idx < len(prevSnap) && idx < len(curSnap) {
					delta := prevSnap[idx].Votes - curSnap[idx].Votes
					if delta > 0 {
						totalDrop += delta
					}
				}
			}
			sumRecipients := 0.0
			for _, v := range cur.SurplusReceived {
				sumRecipients += v
			}
			total := sumRecipients + cur.SurplusExhaustedDelta
			if !floatAlmostEqual(total, totalDrop, tol) {
				t.Fatalf("round %d surplus conservation failed: got %.2f (recipients %.2f + exhausted %.2f), want %.2f (drop in elected)", cur.Round, total, sumRecipients, cur.SurplusExhaustedDelta, totalDrop)
			}
		}
	}
}

func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
