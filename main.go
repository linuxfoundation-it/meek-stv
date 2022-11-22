package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/blackgreen100/meek-stv/election"
	"github.com/blackgreen100/meek-stv/meekstv"
)

func main() {
	f, err := os.Open("election13.txt")
	if err != nil {
		panic(err)
	}

	data := election.Read(f)
	elected := meekstv.Count(data)

	sort.Slice(elected, func(i, j int) bool {
		return elected[i].Votes >= elected[j].Votes
	})

	fmt.Println("-------------------------")
	fmt.Printf("Results of %s\n", data.Title)
	for _, e := range elected {
		if e.State == meekstv.Elected {
			fmt.Printf("%s is elected with %.02f votes\n", e.Name, e.Votes)
		}
	}
}
