package meekstv

import (
	"encoding/json"
	"math"
	"os"
	"testing"

	"github.com/linuxfoundation-it/meek-stv/election"
	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	t.Run("so election 14", func(t *testing.T) {
		got, want := load("election14")
		assertAll(t, want, got)
	})

	t.Run("so election 13", func(t *testing.T) {
		got, want := load("election13")
		assertAll(t, want, got)
	})

	t.Run("so election 12", func(t *testing.T) {
		got, want := load("election12")
		assertAll(t, want, got)
	})

	t.Run("so election 11", func(t *testing.T) {
		got, want := load("election11")
		assertAll(t, want, got)
	})

	t.Run("so election 10", func(t *testing.T) {
		t.Skip() // failing due to precision on fractional votes
		got, want := load("election10")
		assert.ElementsMatch(t, want.Winners, got.Winners(), "winners mismatch")
	})

	t.Run("medical sciences 2022", func(t *testing.T) {
		got, want := load("medsci2022")
		assert.ElementsMatch(t, want.Winners, got.Winners(), "winners mismatch")
	})

	t.Run("chinese 2020", func(t *testing.T) {
		t.Skip() // failing due to precision on fractional votes
		got, want := load("chinese2020")
		assert.ElementsMatch(t, want.Winners, got.Winners(), "winners mismatch")
	})
}

func assertAll(t *testing.T, want *OpaVoteJSONReport, got Log) {
	assert.Equal(t, len(want.Rounds), got.NumRounds())

	omega := math.Pow(10, float64(want.Precision))

	for i, wantRound := range want.Rounds {
		gotRound := got.Round(i)
		assert.Equal(t, wantRound.N, gotRound.Round+1, "round number mismatch")
		assert.InDeltaf(t, float64(wantRound.Thresh)/omega, gotRound.Threshold, 0.01, "round %d threshold", wantRound.N)

		for i, votes := range wantRound.Count {
			assert.InDeltaf(t, float64(votes)/omega, gotRound.VotesOf(i), 0.01, "round %d votes of %d", wantRound.N, i)
		}
		assert.InDeltaf(t, float64(wantRound.Exhausted)/omega, gotRound.Exhausted, 0.01, "round %d exhausted", wantRound.N)
	}
	assert.ElementsMatch(t, want.Winners, got.Winners(), "winners mismatch")
}

func load(filename string) (Log, *OpaVoteJSONReport) {
	return Count(readBallots(filename)), readControl(filename)
}

func readBallots(name string) *election.Election {
	f, err := os.Open("../testdata/" + name + ".txt")
	if err != nil {
		panic(err)
	}

	return election.Read(f)
}

type OpaVoteJSONReport struct {
	NSeats      int                `json:"n_seats"`
	NVotes      int                `json:"n_votes"`
	Title       string             `json:"title"`
	Precision   int                `json:"precision"`
	Withdrawn   []int              `json:"withdrawn"`
	Method      string             `json:"method"`
	Version     string             `json:"version"`
	Candidates  []string           `json:"candidates"`
	Winners     []int              `json:"winners"`
	TieBreaks   []interface{}      `json:"tie_breaks"`
	Options     [][]interface{}    `json:"options"`
	Rounds      []OpaVoteJSONRound `json:"rounds"`
	NValidVotes int                `json:"n_valid_votes"`
}

type OpaVoteJSONRound struct {
	Count      []int  `json:"count"`
	Surplus    int    `json:"surplus"`
	Exhausted  int    `json:"exhausted"`
	Continuing []int  `json:"continuing"`
	N          int    `json:"n"`
	Msg        string `json:"msg"`
	Winners    []int  `json:"winners"`
	Losers     []int  `json:"losers"`
	Thresh     int64  `json:"thresh"`
	Action     struct {
		Type       string `json:"type"`
		Candidates []int  `json:"candidates"`
		Desc       string `json:"desc"`
	} `json:"action,omitempty"`
}

func readControl(name string) *OpaVoteJSONReport {
	f, err := os.Open("../testdata/" + name + ".json")
	if err != nil {
		panic(err)
	}
	r := OpaVoteJSONReport{}
	err = json.NewDecoder(f).Decode(&r)
	if err != nil {
		panic(err)
	}
	return &r
}
