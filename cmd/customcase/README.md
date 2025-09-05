# Custom Case Runner

A small, focused runner to reproduce and debug specific elections with custom ballots. Useful for verifying transfer attribution, exhausted deltas, and round-by-round outcomes when integrating with an external voting service.

## What it does

- Builds an in-memory `Election` with your candidates and ballots
- Runs the Meek STV count (`meekstv.Count`)
- Prints round-by-round summaries including:
  - Threshold (absolute and percent)
  - Exhausted votes
  - Candidate keep factors and votes
  - Surplus transfers and elimination transfers (per-recipient and exhausted deltas)
  - Elected and eliminated candidates per round

## Run

```bash
# from repo root
go run ./cmd/customcase
```

## Plugging in ballots from your voting service

You can modify `cmd/customcase/main.go` to fetch or inject ballots at runtime. Two simple approaches:

1) Hardcode for quick repro (current pattern)
- Set `choices` to your choice IDs in the exact order that defines candidate indices (0..n-1)
- Convert your ballots to zero-based candidate indices and append to the in-memory `Election`

2) Environment or JSON input
- Read `choices` and `ballots` via env vars or a JSON file
- Map your choice IDs to indices using the `choices` array
- Ensure ballots use zero-based indices and weights as integers

Example JSON shape to parse:
```json
{
  "choices": ["choice-id-0", "choice-id-1", "choice-id-2"],
  "seats": 2,
  "ballots": [
    {"weight": 1, "preferences": [2,0,1]}, // means choice-id-2 (rank 1), choice-id-0 (rank 2), choice-id-1 (rank 3)
    {"weight": 1, "preferences": [0,2,1]} // means choice-id-0 (rank 1) ..
  ]
}
```

## Mapping rules (critical)

- Candidate indices are positional: the order of `choices` defines Index 0..n-1
- Each ballot `preferences` array must contain these zero-based indices
- If you use choice IDs elsewhere, always resolve by Index to prevent attribution mix-ups

## Interpreting the output

- Surplus transfers: Effects of an election from the previous round. Shows per-recipient gains
- Elimination transfers: Effects of an elimination from the previous round. Shows per-recipient gains 

## Tips

- Use this runner to iterate quickly on corner cases without changing the main program 