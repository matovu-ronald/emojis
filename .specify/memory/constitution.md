<!-- Sync Impact Report
Version change: N/A → 1.0.0
Modified principles: PRINCIPLE_1_NAME → Test-First Reliability; PRINCIPLE_2_NAME → Stable Search API; PRINCIPLE_3_NAME → Curated Emoji Data; PRINCIPLE_4_NAME → Performance With Simplicity; PRINCIPLE_5_NAME → Usability & Docs
Added sections: Core Principles (materialized), Operational Constraints, Development Workflow & Quality Gates, Governance
Removed sections: None
Templates requiring updates: ✅ .specify/templates/plan-template.md (Constitution Check aligns with principles) | ✅ .specify/templates/spec-template.md (no conflicting guidance) | ✅ .specify/templates/tasks-template.md (task grouping remains valid) | ✅ README.md (runtime guidance referenced)
Follow-up TODOs: None
-->

# emojis Constitution

## Core Principles

### Test-First Reliability
All changes MUST land with automated tests that fail before implementation and pass after. Tests MUST be deterministic, table-driven where practical, and cover new code paths and data mutations in `search`. Regression tests are REQUIRED for any bug fix.

### Stable Search API
Public contracts in the `search` package MUST remain backward compatible within a major version. Breaking changes require a major version bump and a documented migration note. API behavior MUST be deterministic: identical inputs yield identical outputs regardless of environment.

### Curated Emoji Data
Emoji records MUST remain deduplicated, accurately tagged, and accompanied by descriptive labels. Any data addition or update MUST include tests that prove inclusion/exclusion logic and guard against regressions in existing queries. The repository MUST own its emoji data; no runtime network fetches.

### Performance With Simplicity
Keep the search path simple and fast: prefer linear, allocation-light scans over premature indexing. Avoid new dependencies unless they measurably improve correctness or performance. Go code MUST be formatted with `gofmt` and kept idiomatic.

### Usability & Docs
README examples and GoDoc comments MUST reflect current behavior. New exported symbols require concise documentation and, when helpful, example tests. Usage guidance MUST stay minimal and copy-pastable for consumers.

## Operational Constraints

- Target Go toolchain: 1.25.x or higher. Contributors MUST run `gofmt` and `go vet` locally before proposing changes.
- No runtime network or filesystem dependencies for search; the emoji dataset ships with the module.
- Testing MUST run via `go test ./...` and remain under reasonable runtime (<5s on typical laptops) to preserve fast feedback.
- Keep dependency footprint minimal; introduce third-party packages only when justified in the PR description with measurable benefit.
- Any generated data or scripts MUST be reproducible and checked in alongside instructions.

## Development Workflow & Quality Gates

- Every PR MUST state which principle(s) it affects and how compliance was validated.
- Pre-merge gates: `go test ./...`, `gofmt`, and `go vet` MUST pass. Include benchmark evidence when performance-sensitive changes are made.
- Data changes (emoji additions/removals/tag updates) MUST include rationale, coverage by tests, and verification that existing queries remain stable.
- Documentation updates MUST accompany any behavior change and keep README examples in sync.

## Governance

- This constitution supersedes prior informal practices for the repository.
- Amendments require a PR with: the proposed text change, rationale tied to a principle or constraint, and an assessment of ripple effects on templates and docs. Approval requires at least one maintainer review.
- Constitution uses semantic versioning: MAJOR for incompatible governance changes, MINOR for new principles/sections or materially expanded rules, PATCH for clarifications. Version and amendment date MUST update with each accepted change.
- Compliance is reviewed in every PR; non-compliant changes are blocked until aligned. Runtime guidance and examples live in README.md and MUST be kept consistent with this constitution.

**Version**: 1.0.0 | **Ratified**: 2026-01-18 | **Last Amended**: 2026-01-18
