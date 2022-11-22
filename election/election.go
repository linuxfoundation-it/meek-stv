package election

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func Read(in io.ReadCloser) *Election {
	s := bufio.NewScanner(in)
	defer in.Close()

	election := &Election{}

scanner:
	for i := 0; s.Scan(); i++ {
		t := s.Text()

		switch i {
		case 0:
			election.Candidates, election.Seats = parseHeader(t)

		case 1:
			if strings.HasPrefix(t, "-") {
				election.Withdrawn = parseWithdrawn(t)
				continue
			}
			fallthrough

		default:
			if t == "0" {
				break scanner
			}
			election.Ballots = append(election.Ballots, parseBallot(t))
		}
	}

	names := make([]string, 0)
	for s.Scan() {
		names = append(names, s.Text())
	}
	election.CandidateNames = names[:len(names)-1]
	election.Title = names[len(names)-1]

	return election
}

func parseHeader(line string) (candidates, seats int) {
	ss := strings.Split(line, " ")
	is := asInts(ss)
	candidates, seats = is[0], is[1]
	return
}

func parseWithdrawn(line string) map[int]bool {
	ss := strings.Split(line, " ")
	is := asInts(ss)

	out := make(map[int]bool)
	for _, i := range is {
		// substract 1 to make it 0-indexed
		out[i*(-1)-1] = true
	}
	return out
}

func parseBallot(line string) Ballot {
	k := strings.Split(line, " ")
	return Ballot{
		Weight:      asInt(k[0]),
		Preferences: normalize(asInts(k[1 : len(k)-1])),
	}
}

func asInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func asInts(ss []string) (out []int) {
	for _, s := range ss {
		n := asInt(strings.TrimSpace(s))
		if n == 0 {
			continue
		}
		out = append(out, n)
	}
	return
}

// subtract 1 to all candidate indices to make it 0-based
func normalize(ns []int) []int {
	for i := range ns {
		ns[i] = ns[i] - 1
	}
	return ns
}
