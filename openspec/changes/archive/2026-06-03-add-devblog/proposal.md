# Proposal: Add devblog with Kaa persona

**Change**: add-devblog
**Status**: proposed
**Bean**: [jjay-0mni — at merge write a new section in the devblog](../../../.beans/jjay-0mni--at-merge-write-a-new-section-in-the-devblog.md)

## Why

The development of jjay is an interesting journey — choosing a language, building an orchestrator, hitting bugs, evolving the workflow. This raw material should be captured as it happens, not reconstructed later. An AI-generated devblog narrated by the project mascot Kaa (a Eurasian Jay in a jujutsu gi) makes this transparent and fun.

## What Changes

- Add `blog` artifact to the `spec-driven-with-adr` schema
- Blog drafts live in the change directory at `blog/<slug>.md` during development
- At archive time, sync to `devblog/posts/<date>-<slug>.md`
- Narrator is Kaa — first person, brutaal, bossy, commands the flock
- Each post is standalone, short and punchy, about one change's new capabilities
- Fully auto-generated as raw material for a real blog later
- Create retroactive posts for all work done on 2026-06-02

## Capabilities

### New Capabilities

- `devblog`: Blog artifact in schema, persistent posts in `devblog/posts/`, Kaa persona

## Impact

- Modified: `openspec/schemas/spec-driven-with-adr/schema.yaml` (add blog artifact)
- New: `openspec/schemas/spec-driven-with-adr/templates/blog.md` (template)
- New: `devblog/posts/` directory with retroactive posts
- New: `devblog/README.md` introducing Kaa
