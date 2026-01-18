# Implementation Plan: Like-emoji search

**Branch**: `001-like-emoji-search` | **Date**: 2026-01-18 | **Spec**: [specs/001-like-emoji-search/spec.md](specs/001-like-emoji-search/spec.md)
**Input**: Feature specification from `/specs/001-like-emoji-search/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Add a deterministic like-emoji search that accepts an emoji or label, returns a ranked list of related emojis based on shared labels and tags, supports include/exclude tag filters, and enforces a caller-provided (or default) limit while keeping the existing `search` API stable.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.25.x (module uses 1.25.5)  
**Primary Dependencies**: Standard library only; avoid new third-party deps  
**Storage**: In-memory emoji dataset shipped with module; no external storage  
**Testing**: `go test ./...` with table-driven unit tests in `search`  
**Target Platform**: Go module consumers (CLI/apps/services) on macOS/Linux  
**Project Type**: Single library/package (`search`)  
**Performance Goals**: Like-emoji queries return within 0.5s on a typical laptop; deterministic ordering across runs  
**Constraints**: No runtime network/filesystem fetches; gofmt + go vet required; tests under ~5s; deterministic outputs; minimal allocations  
**Scale/Scope**: Small dataset (emoji set) and single-package change; API backwards compatible within v1

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Test-first reliability: Plan includes table-driven tests for new like-emoji search paths; add regression coverage for similarity, filters, limits.
- Stable search API: Extend `search` package without breaking current functions; new entry point keeps deterministic behavior and docs up to date.
- Curated emoji data: No runtime fetches; reuse existing dataset; any data tweaks require tests that protect existing queries.
- Performance with simplicity: Favor linear scan over complex indexing; avoid new deps; keep gofmt/vet clean.
- Usability & docs: Update README/GoDoc as needed; provide quickstart and example usage; deterministic ordering documented.

## Project Structure

### Documentation (this feature)

```text
specs/001-like-emoji-search/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
search/
├── data.go        # emoji dataset
├── search.go      # existing search APIs
├── search_test.go # tests
└── doc.go         # package docs

specs/001-like-emoji-search/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
└── contracts/
```

**Structure Decision**: Single-package Go module; feature work stays in `search/` with tests in the same package and docs in `doc.go`/README; feature documentation under `specs/001-like-emoji-search/`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
