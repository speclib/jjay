## Context

jjay's development journey is worth capturing — language choices, architecture decisions, bugs encountered. Rather than writing a blog post-hoc, we capture raw material as each change is made. The narrator is Kaa, the project mascot: a Eurasian Jay (vlaamse gaai) in a karate gi, who commands a flock of AI agents.

## Goals / Non-Goals

**Goals:**
- Add blog artifact to the schema (same pattern as ADRs)
- Create retroactive posts for all 2026-06-02 work
- Establish Kaa's voice and the devblog structure

**Non-Goals:**
- Publishing infrastructure (static site, RSS — future)
- Polished prose (this is raw material for a real blog)
- Images or illustrations in posts (just markdown text)

## Decisions

### Blog artifact follows ADR pattern

```
During change:                     At archive time:
change/blog/<slug>.md         →    devblog/posts/<date>-<slug>.md
```

Same lifecycle as ADRs: draft in change dir, sync to persistent location at archive. The blog artifact requires `proposal` (needs to know what the change is about).

_Alternative: generate blog post at archive time only — rejected because the draft should be reviewable during the change._

### Blog artifact is optional

Not every change needs a blog post. The artifact exists in the schema but can be skipped (like ADRs). Tooling changes, config tweaks, etc. don't need a post.

### Kaa's voice

- First person: "I", "my flock", "my agents"
- Brutaal, bossy, confident — he's in charge
- Short sentences, punchy paragraphs
- Focuses on what and why, not implementation details
- Dutch personality, English text — occasional Dutch flavor is fine
- Each post opens with a one-line hook

### Post structure

```markdown
# <Title>

<One-line hook>

<2-4 short paragraphs: what was built, why, what's next>
```

No frontmatter, no metadata. Keep it raw and simple.

### Retroactive posts

Six posts for 2026-06-02 work:
1. Choosing Go
2. Project scaffold
3. Spawn command
4. Workspace staleness bug
5. Release process
6. Cleanup command

## Risks / Trade-offs

- [AI voice may feel repetitive] → Vary sentence structure, keep posts short. This is raw material anyway.
- [Blog posts add noise to changes] → Optional artifact, skip when not needed.
