package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/blackgreen100/meek-stv/election"
	"github.com/blackgreen100/meek-stv/meekstv"
)

func main() {
	f, err := os.Open("testdata/election14.txt")
	if err != nil {
		panic(err)
	}

	data := election.Read(f)
	report := meekstv.Count(data)

	report.Print()

	elected := report.Results()
	sort.Slice(elected, func(i, j int) bool {
		return elected[i].Votes >= elected[j].Votes
	})
	elected = elected[:data.Seats]

	fmt.Println("-------------------------")
	fmt.Printf("Results of %s\n", data.Title)
	for _, e := range elected {
		if e.State == meekstv.Elected {
			fmt.Printf("%s is elected with %.02f votes\n", e.Name, e.Votes)
		}
	}
}
