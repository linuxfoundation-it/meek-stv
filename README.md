![Build Status](https://github.com/blackgreen100/meek-stv/actions/workflows/go.yml/badge.svg)

### Description

This repository contains an implementation of MeekSTV vote counting algorithm, 
slightly adapted from [this resource](https://web.archive.org/web/20210225045400/https://prfound.org/resources/reference/reference-meek-rule/), referenced by [Wikimedia foundation](https://meta.wikimedia.org/wiki/Wikimedia_Foundation_elections/Single_Transferable_Vote).

More information about MeekSTV can be found [on Wikipedia](https://en.wikipedia.org/wiki/Counting_single_transferable_votes#Meek).

This vote counting algorithm is the one currently in use at Stack Overflow for moderator elections. 

The software they use is at [OpaVote](https://www.opavote.com/). Although the OpaVote UI allows unregistered users to recount ballots for any given election, it doesn't offer the possibility to tweak many parameters, notably the number of seats available.

This implementation allows to re-run an election with a different number of seats. Currently, you can do that by manually changing the first line of the ballot file.

The first line of OpaVote ballot files appears in the form
```text
6 2
```
where the first number (in this case `6`) is the number of candidates vying for the position, 
and the second number is the number (in this case `2`) is the number of available seats.

By manually changing this line, you can effectively re-run an election and see how the algorithm would determine winning candidates with a different number of available seats.

### Usage

At this time you can only run this program on your local machine with a Go installation.

1. Clone this repository
2. [Install Go](https://go.dev/doc/install)
3. Run the program with `go run ./...`

The program will print some logs that summarize the vote counting. For example, using Stack Overflow 13th moderator election ballot file, 
it will print:

```text
round 1
threshold 9318.33 (33.33)
eliminating "Daniel Widdis"
-------------------------
round 2
threshold 9220.00 (33.33)
eliminating "Ryan M"
-------------------------
round 3
threshold 9035.00 (33.33)
elected "Zoe" with 9208.00 votes
elected "Stephen Rauch" with 9356.00 votes
-------------------------
-------------------------
Results of "Stack Overflow Moderator Election 2021"
"Stephen Rauch" is elected with 9356.00 votes
"Zoe" is elected with 9208.00 votes
```

### How to change the ballot file

At this time you can only change the file name in the program source.
1. Download the ballot .txt file, add it inside your local copy of this repo 
2. In `main.go` at line 13, change the name of the file with the new one.
3. Re-run the program.

### Differences with OpaVote UI

The OpaVote UI aggregates and displays the results of each counting round in such a way that it results 
simpler to understand for everyone. 

For example, it much more clearly displays transfers of vote surpluses (compare with 12th election ballot file).
I _think_ that my implementation is functionally equivalent to OpaVote's one — except for the missing tie breaking function — even if the round progress is reported differently.
The results should be just the same, anyway. If you find any bugs, please report them in the issue tracker.

### Limitations

This program doesn't use fixed point math and the decimals aren't rounded. In practice, this shouldn't cause any appreciable differences from the official OpaVote counting method, 
however *it just might*, in very tight election rounds. Always refer to the OpaVote counting algorithm for official results. At your discretion, report bugs in the issue tracker.   

### Disclaimer

This software is in no way a replacement for Stack Overflow's own election process. 
It is simply meant to make it easier to recount votes with different parameters. The results of any recount done with this program aren't binding and don't prove or disprove anything. 
This program is exclusively meant to satisfy my — and hopefully your — curiosity. 

### Improvements

- add debug statements
- add tie breaking
- add command line flag parsing, so that the program can be run as a standalone executable
- allow specifing ballot file from command line
