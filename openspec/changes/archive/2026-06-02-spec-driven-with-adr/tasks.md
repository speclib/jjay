# Tasks: spec-driven-with-adr

## 1. Schema Setup

- [x] 1.1 Fork spec-driven schema via `openspec schema fork spec-driven spec-driven-with-adr`
- [x] 1.2 Add `adr` artifact to `openspec/schemas/spec-driven-with-adr/schema.yaml` (requires proposal, generates to `openspec/adrs/`)
- [x] 1.3 Create ADR template at `openspec/schemas/spec-driven-with-adr/templates/adr.md`
- [x] 1.4 Validate the schema with `openspec schema validate spec-driven-with-adr`

## 2. Config Switch

- [x] 2.1 Update `openspec/config.yaml` schema field from `spec-driven` to `spec-driven-with-adr`
- [x] 2.2 Verify with `openspec schemas` that the new schema is listed and active

## 3. Retroactive ADRs

- [x] 3.1 Create `openspec/adrs/` directory
- [x] 3.2 Write ADR-001: Use Go as implementation language (reference archived proposal)
- [x] 3.3 Write ADR-002: OpenSpec config — project context and light rules (reference archived proposal)
