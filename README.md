## meek-stv

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

### TODOs

- add debug statements
- add tie breaking
- add command line flag parsing, so that the program can be run as a standalone executable
