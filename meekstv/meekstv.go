package meekstv

import (
	"fmt"
	"math"

	"example.com/election"
)

func Count(params *election.Election) Candidates {
	// Initialize Election
	getInitialState := func(i int) CandidateState {
		if params.Withdrawn[i] {
			return Withdrawn
		}
		return Hopeful
	}
	getInitialKeepFactor := func(i int) float64 {
		if params.Withdrawn[i] {
			return 0.0
		}
		return 1.0
	}

	// Set each candidate’s state to hopeful or withdrawn.
	// Set each hopeful candidate’s keep factor kf to 1, and each withdrawn candidate’s keep factor to 0.
	cs := make(Candidates, params.Candidates)
	for i := 0; i < params.Candidates; i++ {
		cs[i] = &Candidate{
			Index:      i,
			Name:       params.CandidateNames[i],
			State:      getInitialState(i),
			KeepFactor: getInitialKeepFactor(i),
		}
	}

	// Test count complete. Proceed to step C if all seats are filled,
	// or if the number of elected plus hopeful candidates is less than or equal to the number of seats.

	round := &meekStvRound{
		omega:      1 / 10e6,
		candidates: cs,
	}
	for ; ; round.n++ {
		hopeful := cs.countState(Hopeful)
		elected := cs.countState(Elected)
		if elected >= params.Seats || elected+hopeful <= params.Seats {
			return round.complete(params.Seats)
		}
		fmt.Println("round", round.n+1)
		round.run(params)
		fmt.Println("-------------------------")
		// Continue. Proceed to the next round at step B.1.
	}
}

// holds state of a MeekSTV count
type meekStvRound struct {
	n           int
	candidates  Candidates
	omega       float64
	prevSurplus float64
}

func (round *meekStvRound) run(el *election.Election) {
	newlyElected := false

	round.candidates.resetVotes()
	// iterate
	// Distribute votes.
	// For each ballot: set ballot weight w to 1, and then for each candidate, in order of rank on that ballot:
	// add w multiplied by the keep factor kf of the candidate (to 9 decimal places, rounded up)
	// to that candidate’s vote v, and reduce w by the same amount, until no further candidate remains
	// on the ballot or until the ballot’s weight w is 0.
	for _, bl := range el.Ballots {
		w := float64(0.0)
		for i := 0; i < len(bl.Preferences); i++ {
			c := round.candidates[bl.Preferences[i]]
			w = float64(bl.Weight)
			v := w * c.KeepFactor
			c.Votes += v
			w -= v
			if w <= 0.0 {
				break
			}
		}
	}

	// Update quota. Set quota q to the sum of the vote v for all candidates (step B.2.a),
	// divided by one more than the number of seats to be filled,
	// truncated to 9 decimal places, plus 0.000000001 (1/109).
	elected := round.candidates.countState(Elected)
	totvotes := round.candidates.countVotes()
	threshold := totvotes / (1.0 + float64(el.Seats) - float64(elected))
	fmt.Printf("threshold %.02f (%.02f)\n", threshold, threshold/totvotes*100)

	// Find winners. Elect each hopeful candidate with a vote v greater than or equal to the quota (v ≥ q).
	for _, c := range round.candidates {
		if c.Votes >= threshold {
			c.State = Elected
			// c.KeepFactor = (c.KeepFactor * threshold) / c.Votes
			newlyElected = true
			fmt.Printf("elected %s with %.02f votes\n", c.Name, c.Votes)
		}
	}

	// Update keep factors. Set the keep factor kf of each elected candidate to the candidate’s
	// current keep factor kf, multiplied by the current quota q (to 9 decimal places, rounded up),
	// and then divided by the candidate’s current vote v (to 9 decimal places, rounded up).
	// Continue iteration at step B.2.a.
	for _, c := range round.candidates {
		if c.State == Elected {
			c.KeepFactor = (c.KeepFactor * threshold) / c.Votes
		}
	}

	// Calculate the total surplus s, as the sum of the individual surpluses (v – q) of the elected candidates,
	// but not less than 0.
	totSurplus := 0.0
	for _, c := range round.candidates {
		c.Surplus = math.Max(c.Votes-threshold, 0.0)
		totSurplus += c.Surplus
	}

	// Test for iteration finished. If step B.2.c elected a candidate, continue at B.1.
	if newlyElected {
		round.prevSurplus = totSurplus
		return
	}

	// Otherwise, if the total surplus s is less than omega, or (except for the first iteration)
	// if the total surplus s is greater than or equal to the surplus s in the previous iteration, continue at B.3.
	if totSurplus < round.omega || (round.n > 0 && totSurplus >= round.prevSurplus) {
		goto defeatCandidates
	}

	// Defeat low candidate.
	// Defeat the hopeful candidate c with the lowest vote v, breaking any tie per procedure T,
	// where each candidate c' is tied with c if vote v' for c' is less than or equal to v plus total surplus s.
	// Set the keep factor kf of c to 0.
defeatCandidates:
	var d = &Candidate{Votes: math.MaxFloat64}
	for _, c := range round.candidates {
		if c.State == Hopeful && c.Votes < d.Votes {
			d = c
		}
	}
	d.State = Defeated
	d.KeepFactor = 0.0
	fmt.Println("eliminating", d.Name)

	// Continue. Proceed to the next round at step B.1.
	round.prevSurplus = totSurplus
}

// tiebreaking
// Ties can arise in B.3, when selecting a candidate for defeat.
// Use the defined tiebreaking procedure to select for defeat one candidate from the group of tied candidates.

func (round *meekStvRound) complete(seats int) Candidates {
	// Elect remaining. If any seats are unfilled, elect remaining hopeful candidates.
	elected := round.candidates.countState(Elected)
	if elected < seats {
		for _, c := range round.candidates {
			if c.State == Hopeful {
				c.State = Elected
			}
		}
	}

	// Defeat remaining. Otherwise defeat remaining hopeful candidates.
	for _, c := range round.candidates {
		if c.State != Elected {
			c.State = Defeated
		}
	}

	return round.candidates
}
