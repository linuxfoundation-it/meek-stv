package main

import (
	"fmt"
	"os"
	"sort"

	"example.com/election"
	"example.com/meekstv"
)

func main() {
	f, err := os.Open("election12.txt")
	if err != nil {
		panic(err)
	}

	elected := meekstv.Count(election.Read(f))

	sort.Slice(elected, func(i, j int) bool {
		return elected[i].Votes >= elected[j].Votes
	})

	fmt.Println("-------------------------")
	for _, e := range elected {
		if e.State == meekstv.Elected {
			fmt.Printf("%s is elected with %.02f votes\n", e.Name, e.Votes)
		}
	}
}
