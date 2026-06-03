# Proposal: GitHub Pages website with devblog

**Change**: add-pages-website
**Bean**: [jjay-4grr — create github pages website with sub for devlog](../../../.beans/jjay-4grr--create-github-pages-website-with-sub-for-devlog.md)

## Why

jjay has a mascot (Kaa, the blue jay in a gi), a hero banner, and a 14-post devblog written in Kaa's voice — but no public home. The devblog is currently just markdown files in `devblog/posts/`, invisible to anyone not reading the repo. A compact website gives the project a front door for advanced developers and turns the devblog into something people actually read, without changing how posts are authored.

## What Changes

- Add a **Hugo** site under `site/` that builds to `https://speclib.github.io/jjay/` (project page on a subpath).
- **Frontpage** (`/`): compact, advanced-dev framing — logo, "Control the flock." one-liner, a short what-it-is, the single `jjay spawn` command that replaces the manual dance, a 4-item feature grid (spawn / merge / cleanup / session-open), and links to GitHub + the devblog.
- **Devblog** (`/blog/`): the existing `devblog/posts/` directory **mounted directly** as Hugo content via a module mount — no copy, no transform step. Posts stay raw (no frontmatter), exactly as Kaa writes them. Date comes from the `YYYY-MM-DD-` filename prefix; the post title is extracted from the first `# ` heading (single source of truth).
- **Style**: terminal/hacker aesthetic — dark, monospace, fast, keyboard-friendly — with the mascot's blue as the accent and `jjay-persona.png` as the logo. Not boring, but respects the advanced-dev audience's taste.
- **GitHub Action** (`.github/workflows/pages.yml`): build with Hugo and deploy via the GitHub Actions → Pages artifact path (`actions/deploy-pages@v4`), triggered on push to `main` touching `site/**` or `devblog/posts/**`. Archiving a change (which appends a post) therefore auto-updates the site.

## Capabilities

### New Capabilities

- `pages-website`: A Hugo-built GitHub Pages site with a compact frontpage and a devblog that reads `devblog/posts/` directly, deployed automatically via GitHub Actions.

### Modified Capabilities

_(none — this is additive; no existing specs change)_

## Impact

- **New**: `site/` — Hugo site (`hugo.toml`, `content/_index.md`, `layouts/`, `assets/css/`, `static/`).
- **New**: `.github/workflows/pages.yml` — build + deploy workflow.
- **New**: copies of `jjay-persona.png` / `hero.png` into `site/static/` (source remains `artwork/`).
- **Unchanged**: `devblog/posts/*.md` — mounted, never modified. Kaa's no-frontmatter rule stays intact.
- **One-time repo setting**: Settings → Pages → Source = "GitHub Actions". Org-level (`speclib`) must allow Pages + Actions deploy.
- **No Go code changes**, no new Go dependencies. Hugo is a build-time tool, pinned in the Action.

## Non-goals

- No custom domain (stays on `speclib.github.io/jjay/`, no CNAME).
- No search, comments, or analytics — advanced-dev audience doesn't want them and the bean doesn't ask.
- No reusable/packaged Hugo theme — bespoke `layouts/` is simpler.
- No frontmatter added to existing posts — the raw-post convention is a brand decision and is preserved.
- RSS is out of scope as a task (Hugo offers it nearly free; noted as a future nice-to-have, not built here).
