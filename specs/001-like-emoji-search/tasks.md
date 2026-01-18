# Tasks: Like-emoji search

**Input**: Design documents from `/specs/001-like-emoji-search/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are required per constitution; write failing tests before implementation within each story.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Ensure tooling and baseline cleanliness for Go module and tests.

- [x] T001 Run gofmt on codebase and verify clean working tree in search/ and root
- [x] T002 Run go vet for package search/ to confirm no baseline issues
- [x] T003 [P] Add constitution compliance note for like-emoji work in README section if needed in README.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core shared components required before user stories.

- [x] T004 Define LikeParams struct and validation scaffolding in search/search.go (seed fields, include/exclude, limit default)
- [x] T005 [P] Add normalization helpers for labels/tags (trim, lowercase, split terms) in search/search.go
- [x] T006 [P] Add helper to compute overlap score and deterministic sort comparator in search/search.go
- [x] T007 Add unit tests for validation and helpers (normalization, overlap, clamping) in search/search_test.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel.

---

## Phase 3: User Story 1 - Find similar emojis by example (Priority: P1) 🎯 MVP

**Goal**: Provide like-emoji search that returns related emojis for a seed emoji/label with deterministic ordering.

**Independent Test**: Call like search with a single emoji or label and verify related results are returned in stable order without other features.

### Tests for User Story 1 (write first)

- [x] T008 [P] [US1] Add table-driven tests for like search similarity (emoji seed vs label seed) ensuring seed is excluded and order is deterministic in search/search_test.go
- [x] T009 [P] [US1] Add tests for empty/unknown seed returning empty slice in search/search_test.go
- [x] T024 [P] [US1] Add tests ensuring representative seeds return at least three related emojis where applicable (SC-001) in search/search_test.go
- [x] T025 [P] [US1] Add repeatability test running the same query multiple times to confirm identical ordering/content (SC-002) in search/search_test.go

### Implementation for User Story 1

- [x] T010 [US1] Implement Like function core: seed resolution, similarity scoring, deterministic ordering in search/search.go
- [x] T011 [US1] Add GoDoc comment and example for Like in search/doc.go
- [x] T012 [US1] Update README usage section with basic like search example in README.md

**Checkpoint**: User Story 1 fully functional and independently testable.

---

## Phase 4: User Story 2 - Filter similar results by tags (Priority: P2)

**Goal**: Allow include/exclude tag filters to refine similar emoji results.

**Independent Test**: Execute like search with include/exclude tags and verify results respect filters without needing other stories.

### Tests for User Story 2 (write first)

- [x] T013 [P] [US2] Add tests covering include-only, exclude-only, combined filters, and precedence (exclude wins) in search/search_test.go
- [x] T026 [P] [US2] Add include/exclude compliance matrix tests across at least five tag combinations to verify filter correctness (SC-003) in search/search_test.go

### Implementation for User Story 2

- [x] T014 [US2] Implement include/exclude filtering path before scoring, honoring exclude precedence in search/search.go
- [x] T015 [US2] Extend README filter example and quickstart to show include/exclude usage in README.md and specs/001-like-emoji-search/quickstart.md

**Checkpoint**: User Stories 1 and 2 independently functional and filter logic validated.

---

## Phase 5: User Story 3 - Control limits and ordering (Priority: P3)

**Goal**: Enforce caller-provided/default limits and keep deterministic ordering with tie-breakers.

**Independent Test**: Run like search with limit; verify capped length and stable ordering across runs.

### Tests for User Story 3 (write first)

- [x] T016 [P] [US3] Add tests for limit clamping (negative → 0, large → dataset size/default 50) and ordering stability across repeated calls in search/search_test.go

### Implementation for User Story 3

- [x] T017 [US3] Apply limit enforcement and early-stop collection; confirm tie-breaker (label asc) in search/search.go
- [x] T018 [US3] Document limit and ordering rules in doc.go and README.md

**Checkpoint**: All user stories independently functional with deterministic capped results.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final quality, documentation, and verification.

- [x] T019 [P] Add/refresh package-level comments or examples in search/doc.go to align with new Like API
- [x] T020 [P] Verify quickstart compiles/runs and matches API behavior in specs/001-like-emoji-search/quickstart.md
- [x] T021 [P] Run full test suite and vet (go test ./..., go vet ./...) before release
- [x] T022 Review for deterministic outputs and add any missing regression cases in search/search_test.go
- [x] T023 Update spec and plan statuses if needed to reflect completion in specs/001-like-emoji-search/
- [x] T027 [P] Add benchmark or timed test to confirm like search completes within ~0.5s on full dataset (SC-004); document result in README or quickstart

---

## Dependencies & Execution Order

- Phase order: Setup → Foundational → US1 → US2 → US3 → Polish.
- User Story dependencies: US1 is independent once Foundational is done; US2 depends on US1 implementation (shared Like function) but tests can be authored earlier; US3 depends on US1 core and may share filtering code from US2.

## Parallel Execution Examples

- Foundational: T005 and T006 can run in parallel after T004; T007 follows to cover helpers.
- US1: T008 and T009 can run in parallel; T011/T012 can run in parallel after T010.
- US2: T013 can start early; T015 can run in parallel with T014 once API shape is stable.
- US3: T016 can start after T010; T018 can run in parallel with T017 once limit behavior is defined.
- Polish: T019–T022 can run in parallel as they touch different files, with T021 last.

## Implementation Strategy

- MVP first: Complete Foundational → US1, validate with tests before proceeding.
- Incremental: Layer filters (US2) then limits/ordering (US3), each with dedicated tests.
- Maintain determinism: Keep scoring/sort pure and documented; add regression tests when modifying data or rules.
