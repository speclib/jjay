## Beans OpenSpec Integration

When I refer to issues like <this-dirname>-rn3b checkout the task
in @.beans/<this-dirname>-rn3b-*.md

In this project we will use these tasks as epics for making openspec proposals.

Run the `beans prime` command to know how beans work.

WHEN you create a proposal at a link to this task in the proposal.md.
WHEN a bean is used to create an proposal change the status to "in-progress"
WHEN a proposal is archived add the link to the archived proposal in the frontmatter of this task like this:

```
openspec-link: openspec/changes/archive/....
```

You are allowed to update these statuses in the task frontmatter:

- in-progress
- todo
- draft
- completed
- scrapped

When making changes you are allowed to update the date/time in `updated_at` in the task frontmatter

Besides updating status and openspec-link, you are NOT ALLOWED to modify the contents of the task file.

## Git/jj

- We use jj
- I don't want Claude to appear as committer or coworker in the git-history

## OpenSpec Archive triggers

WHEN /opsx:archive is called the work should be made ready for merging in bookmark main.

Run these steps in order:

1. **Blog gate (auto-create).** Before anything else, check `openspec status --change "<name>" --json`. If the `blog` artifact is not `done`:
   - Run `openspec instructions blog --change "<name>" --json` to get the template and context.
   - Write the blog **retrospectively**: read `proposal.md` and the completed `tasks.md` for context, and narrate what was *actually* built (the outcome), not just what was planned.
   - Re-check `openspec status` to confirm `blog` is now `done` before proceeding.
   - Only `blog` is auto-created. Other incomplete artifacts (e.g. `adr`) remain **warn-only** — surface them as a warning and let the user confirm proceeding, do not auto-create them.
2. jj describe
3. the changelog should also be updated. Each entry headline (the `- **name** — summary` bullet) MUST be <= 80 characters; put detail and `See [proposal]/[ADR]` links on indented sub-bullets below it, not in the headline.
