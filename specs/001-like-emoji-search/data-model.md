# Data Model: Like-emoji search

## Entities

### Emoji
- **Fields**:
  - `Label` (string, required): Human-readable name; lowercased for matching.
  - `Char` (string, optional): Emoji rune used for seed lookup and output.
  - `Tags` ([]string, required): Descriptive tags; lowercased for matching.
- **Rules**: Tags and label terms drive similarity scoring; no duplicates per emoji; Char may be empty for aliases.

### LikeParams
- **Fields**:
  - `SeedEmoji` (string, optional): Single-rune emoji input.
  - `SeedLabel` (string, optional): Label text input; used when emoji not provided/found.
  - `Include` ([]string, optional): Tags/terms to require (any-match).
  - `Exclude` ([]string, optional): Tags/terms to omit (any-match).
  - `Limit` (int, optional): Max results; default 50; clamp to [0, dataset size].
- **Rules**:
  - At least one of `SeedEmoji` or `SeedLabel` must be non-empty after trimming.
  - Normalize all terms to lowercase and trim whitespace.
  - Exclude filter takes precedence over include.

### LikeEmojiResult
- **Fields**:
  - `Emojis` ([]Emoji): Ordered by score desc, label asc; excludes the seed emoji.
- **Rules**: Deterministic ordering; length respects `Limit`.

## Relationships & Processing
- Seed lookup: If `SeedEmoji` matches dataset `Char`, use that record; otherwise match label terms.
- Similarity: Based on overlap of normalized label terms + tags between seed and candidates.
- Filtering: Apply exclude filter first, then include filter, then scoring and ordering.
- Limiting: Stop after collecting `Limit` items.
