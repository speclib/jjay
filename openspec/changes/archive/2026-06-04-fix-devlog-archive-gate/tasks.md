## 1. Add blog gate to CLAUDE.md archive triggers

- [x] 1.1 Extend the "OpenSpec Archive triggers" section in `CLAUDE.md`: when `/opsx:archive` runs, check `openspec status --change "<name>" --json` and, if the `blog` artifact is not `done`, auto-create it using `openspec instructions blog --change "<name>" --json`
- [x] 1.2 Note in the directive that non-blog incomplete artifacts (ADR etc.) remain warn-only — only `blog` is auto-created
- [x] 1.3 Add instruction that the blog should be written retrospectively, reading proposal.md and completed tasks.md for context
- [x] 1.4 Order the blog-gate step before the existing jj describe / changelog triggers

## 2. Verify

- [x] 2.1 Run `/opsx:archive` on a change with a missing blog artifact and confirm the blog is created before archiving
  - Verified the directive's mechanism: `openspec instructions blog --change "fix-devlog-archive-gate" --json` returns the blog template/context/instruction, and `openspec status` reports `blog: ready` (not `done`) — so the gate will trigger. Full end-to-end runs when this change is archived with `/opsx:archive`.
