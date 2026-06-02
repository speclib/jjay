## Context

OpenSpec 1.2.0 ships only the `spec-driven` schema. We need a project-local fork that adds an ADR artifact. The `openspec schema fork` command copies the schema into the project for customization.

ADRs record architectural decisions with context, alternatives, and consequences. They persist in `openspec/adrs/` — outside the change lifecycle — so future proposals can reference prior reasoning.

## Goals / Non-Goals

**Goals:**
- Add an `adr` artifact to the schema that generates to `openspec/adrs/`
- Keep the existing artifact flow intact (proposal → specs → design → tasks)
- Create retroactive ADRs for decisions already made (Go language, config)

**Non-Goals:**
- Not building tooling to auto-number ADRs (manual numbering is fine)
- Not modifying the archive behavior (ADRs simply aren't in the change dir)
- Not making ADRs mandatory — the artifact exists but can be skipped

## Decisions

**Fork via `openspec schema fork`**
Use the built-in fork command rather than copying files manually. This creates `openspec/schemas/spec-driven-with-adr/` with the correct structure.
_Alternative: `schema init` from scratch — rejected because we'd lose the existing templates._

**ADR artifact placement in dependency graph**
The ADR artifact requires `proposal` (needs context of what's being decided) but does NOT block `tasks`. It sits alongside `design` — both inform implementation but neither depends on the other.
_Alternative: ADR requires design — rejected because some decisions are made before/without a design._

**ADR output path: `openspec/adrs/<number>-<slug>.md`**
ADRs generate outside the change directory into a shared persistent location. This is the key difference from other artifacts — they survive archival by never being in the change dir.
_Alternative: generate in change dir and copy on archive — rejected because it requires custom archive logic._

**ADR numbering: manual, zero-padded three-digit**
Simple `001-use-go.md` convention. The AI creating the ADR scans existing files to determine the next number.
_Alternative: date-based naming — rejected because sequential numbers show decision order more clearly._

## Risks / Trade-offs

- [Schema fork diverges from upstream] → Acceptable for a project-local schema. If openspec adds `spec-driven-with-adr` upstream, we can switch back.
- [ADR outside change dir is unconventional for openspec] → The `generates` path can point outside the change directory. This is supported but unusual.
- [Retroactive ADRs lack full context] → We'll mark them as retroactive. The archived proposals contain the original reasoning.
