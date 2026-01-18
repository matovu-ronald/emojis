# Contract: Like-emoji search (Go package `search`)

## Function Signature (proposed)

```go
// Like returns emojis similar to a seed emoji or label.
func Like(q LikeParams) []Emoji
```

### Input `LikeParams`
- `SeedEmoji` (string, optional): A single emoji rune. If found in dataset, used as the seed record.
- `SeedLabel` (string, optional): Label text fallback when emoji not provided or not found.
- `Include` ([]string, optional): Tags/terms to require (any-match). Empty means no include filtering.
- `Exclude` ([]string, optional): Tags/terms to omit (any-match). Exclude takes precedence.
- `Limit` (int, optional): Maximum results; default 50; clamp to [0, dataset size].

**Validation**
- Reject (return empty) when both `SeedEmoji` and `SeedLabel` are empty after trimming.
- Normalize all terms: trim, lowercase. Ignore empty terms after normalization.

### Output
- `[]Emoji` ordered deterministically: score desc (overlap of tags + label terms), then label asc.
- Excludes the seed emoji itself from results.
- Length respects `Limit`.

### Determinism & Stability
- Same inputs → identical outputs across runs.
- No external data sources; dataset is in-repo.
- Tie-breaker: label ascending when scores match.

### Error Handling
- Invalid/empty input: return empty slice (no panic).
- Limit < 0 treated as 0; limits above dataset size clamp to dataset size.

### Examples

```go
results := search.Like(search.LikeParams{
    SeedEmoji: "👍",
    Include:   []string{"hand"},
    Exclude:   []string{"skin-tone"},
    Limit:     5,
})

results := search.Like(search.LikeParams{
    SeedLabel: "smile",
    Limit:     10,
})
```
