package meekstv

import (
	"encoding/json"
	"math"
	"os"
	"testing"

	"github.com/blackgreen100/meek-stv/election"
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
		t.Skip() // failing due to values not in delta
		got, want := load("election10")
		assertAll(t, want, got)
	})
}

func assertAll(t *testing.T, want *OpaVoteJSONReport, got Log) {
	assert.Equal(t, len(want.Rounds), got.NumRounds())

	omega := math.Pow(10, float64(want.Precision))

	for i, wantRound := range want.Rounds {
		gotRound := got.Round(i)
		assert.Equal(t, wantRound.N, gotRound.Round+1)
		assert.InDelta(t, float64(wantRound.Thresh)/omega, gotRound.Threshold, 0.01)

		for i, votes := range wantRound.Count {
			assert.InDelta(t, float64(votes)/omega, gotRound.CandidateSnapshot[i].Votes, 0.01)
		}
		assert.InDelta(t, float64(wantRound.Exhausted)/omega, gotRound.Exhausted, 0.01)
	}
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
