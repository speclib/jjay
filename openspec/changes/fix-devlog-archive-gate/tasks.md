## 1. Add blog gate to CLAUDE.md archive triggers

- [ ] 1.1 Extend the "OpenSpec Archive triggers" section in `CLAUDE.md`: when `/opsx:archive` runs, check `openspec status --change "<name>" --json` and, if the `blog` artifact is not `done`, auto-create it using `openspec instructions blog --change "<name>" --json`
- [ ] 1.2 Note in the directive that non-blog incomplete artifacts (ADR etc.) remain warn-only — only `blog` is auto-created
- [ ] 1.3 Add instruction that the blog should be written retrospectively, reading proposal.md and completed tasks.md for context
- [ ] 1.4 Order the blog-gate step before the existing jj describe / changelog triggers

## 2. Verify

- [ ] 2.1 Run `/opsx:archive` on a change with a missing blog artifact and confirm the blog is created before archiving
