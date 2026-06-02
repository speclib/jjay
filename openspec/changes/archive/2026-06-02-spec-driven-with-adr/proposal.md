# Proposal: Switch to spec-driven-with-adr schema

**Change**: spec-driven-with-adr
**Status**: proposed
**Bean**: [jjay-x7qa — switch to spec-driven-with-adr](../../../.beans/jjay-x7qa--switch-to-spec-driven-with-adr.md)

## Why

Architectural decisions get buried in the archive. We already hit this — the "use Go" decision lives in `openspec/changes/archive/2026-06-02-techstack-go/proposal.md` where future contributors won't find it.

ADRs (Architectural Decision Records) persist outside the change lifecycle so reasoning stays discoverable. See: https://intent-driven.dev/blog/2026/04/29/spec-driven-development-with-adr/

## What Changes

- Fork `spec-driven` schema to `spec-driven-with-adr` (project-local)
- Add an `adr` artifact that generates to `openspec/adrs/<number>-<slug>.md`
- The ADR artifact sits alongside design (not archived with the change)
- Switch `openspec/config.yaml` to use the new schema
- Create retroactive ADRs for existing decisions:
  - ADR-001: Use Go as implementation language
  - ADR-002: OpenSpec config — project context and light rules

## Capabilities

### New Capabilities

- `adr-persistence` — ADRs persist in `openspec/adrs/` outside the change lifecycle

### Modified Capabilities

_(none — this is additive to the schema, existing workflow unchanged)_

## Impact

- New directory: `openspec/schemas/spec-driven-with-adr/` (project-local schema)
- New directory: `openspec/adrs/` (persistent ADR storage)
- `openspec/config.yaml` schema field changes from `spec-driven` to `spec-driven-with-adr`

## Non-goals

- Not changing the existing artifact flow (proposal → specs → design → tasks)
- Not making ADRs mandatory for every change — they're for architectural decisions
