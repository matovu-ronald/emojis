# Research: Like-emoji search

## Decisions

### 1) Similarity scoring and ordering
- **Decision**: Use overlap-based scoring on normalized tags + label terms (case-insensitive). Compute score as shared_terms / total_unique_terms (Jaccard-like). Exclude the seed emoji from candidates. Stable sort by score desc, then label asc for determinism.
- **Rationale**: Keeps implementation simple, deterministic, dependency-free, and aligns with existing label/tag data. Sorting rule satisfies constitution requirement for deterministic outputs.
- **Alternatives considered**: Vector/embedding similarity (adds deps and non-determinism), Levenshtein on labels only (ignores tags), popularity/usage weighting (data not available, would add bias and instability).

### 2) Input normalization and seed handling
- **Decision**: Accept either an emoji rune or label text. Normalize by trimming whitespace, lowercasing, and splitting labels/tags on spaces. If the seed is a rune present in the dataset, map it to its record; otherwise treat input as label text for matching. Empty/invalid input yields an empty result without errors.
- **Rationale**: Supports both UI copy/paste of emojis and programmatic label inputs; avoids errors for unknowns while staying deterministic.
- **Alternatives considered**: Require emoji-only input (hurts usability), auto-fuzzy matching across all fields (less deterministic, higher runtime).

### 3) Filtering and limits
- **Decision**: Apply exclude tags first (hard filter), then include tags to keep only items matching at least one include term when provided. Enforce a default limit of 50; clamp limits to [0, dataset_size] and stop once the cap is reached.
- **Rationale**: Exclude precedence matches spec edge cases; default 50 covers typical UI lists; early stop keeps runtime low.
- **Alternatives considered**: Include precedence (would violate spec edge case), unlimited results (hurts performance and UX), default limit of 10 (too restrictive for discovery).

### 4) Determinism guardrails
- **Decision**: All operations use fixed ordering rules (score desc, then label asc). No randomization or time-dependent data. Tests will assert identical outputs across repeated runs for the same inputs.
- **Rationale**: Satisfies constitution (Stable Search API, Performance with Simplicity) and spec acceptance scenarios for repeatable ordering.
- **Alternatives considered**: Secondary sorting by length/popularity (needs extra data, could change over time).
