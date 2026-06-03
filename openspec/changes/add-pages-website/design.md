## Context

jjay has strong brand assets (the Kaa mascot, hero banner, a 14-post devblog in Kaa's brutaal first-person voice) but no public site. The devblog lives as raw markdown in `devblog/posts/` — files named `YYYY-MM-DD-slug.md`, each starting with a `# Title` line and deliberately carrying **no frontmatter** (Kaa's rule: "no frontmatter, no metadata, keep it raw"). Posts are appended at archive time by `/opsx:archive`.

The audience is advanced developers who use advanced tools (jj, tmux, agents). That audience reads "not boring" as *fast, tasteful, doesn't waste my time* — not animations and gradients. The design target is a terminal/hacker aesthetic (dark, monospace) with the mascot supplying the personality.

## Goals / Non-Goals

**Goals:**
- A compact frontpage and a `/blog/` devblog at `https://speclib.github.io/jjay/`.
- The devblog builds from `devblog/posts/` **directly**, with no copy or transform — that directory stays the single source of truth and stays raw.
- Auto-deploy: pushing a new/changed post (e.g. via archive) rebuilds and publishes the site.
- Mascot-forward but credible for the advanced-dev audience.

**Non-Goals:**
- Custom domain, search, comments, analytics, packaged theme, RSS-as-a-task.
- Modifying existing posts or adding frontmatter to them.
- Any Go code or new Go dependencies.

## Decisions

### Hugo, with `devblog/posts/` mounted (no transform)
Hugo's module mounts let us mount `../devblog/posts` as `content/blog` with zero copying:

```toml
[[module.mounts]]
  source = "../devblog/posts"
  target = "content/blog"
```

The site (`site/`) owns layout and style only; `devblog/posts/` owns content. Archive a change → post appears → site rebuilds. No generated directory to gitignore, no duplication. This matches jjay's lean, no-duplication ethos.

### Post title = the H1, extracted at render time (Path 2)
Hugo parses the `YYYY-MM-DD-` filename prefix into `.Date` natively, so date and ordering are free. The **title** is not free — Hugo wants it in frontmatter, which the posts don't have. Rather than relax the no-frontmatter rule, the title is extracted from the first `# ` line of `.RawContent` inside a contained partial in `site/layouts/` (used by both the list and single templates). The posts never learn the site exists.

This is the single real technical risk in the project. It is contained to `layouts/partials/` and has a defined fallback (see risk ledger). The H1 *is* the canonical title by Kaa's own format, so this keeps one source of truth rather than introducing a competing filename-derived title.

### Project page on a subpath → relative URLs everywhere
`baseURL = "https://speclib.github.io/jjay/"`. Because the site is served under `/jjay/`, **every** asset and link reference must go through Hugo's `relURL` / `.RelPermalink` / `|absURL` — never a leading-slash literal like `/css/main.css`, which works locally but 404s on Pages. This is enforced by convention in the layouts and is a verification step in tasks.

### Deploy via GitHub Actions → Pages artifact
Use the modern path (`actions/configure-pages`, `actions/upload-pages-artifact`, `actions/deploy-pages@v4`) rather than pushing a `gh-pages` branch. No extra branch, clean history, official Hugo-on-Pages pattern. Workflow needs `pages: write` and `id-token: write` permissions. Hugo version is pinned in the workflow to avoid drift. Trigger: push to `main` with `paths: [site/**, devblog/posts/**]`.

### Frontpage: pitch + magic command + 4-item feature grid
Compact enough to grasp in one screen: logo (`jjay-persona.png`), "Control the flock.", a 2–3 line what-it-is, one code block showing the `jjay spawn` command that replaces the whole manual workflow, a 4-item grid (spawn / merge / cleanup / session-open), and links to GitHub + devblog. No README dump, no scroll-to-understand.

## Data flow (the living-artifact loop)

```
/opsx:archive a change
   │  writes devblog/posts/<date>-<slug>.md  (Kaa's voice, raw)
   ▼
git commit + push to main
   ▼
GH Action (paths: site/**, devblog/posts/**)
   │  hugo --baseURL https://speclib.github.io/jjay/
   │  mounts ../devblog/posts → content/blog
   │  layouts: .Date from filename, title from H1 partial
   ▼
actions/deploy-pages@v4  →  site live, blog auto-updated
```

The devblog becomes a side effect of the archive workflow — no separate "publish" step.

## File layout

```
site/
├── hugo.toml                  baseURL .../jjay/ ; module.mount ../devblog/posts → content/blog
├── content/_index.md          frontpage copy (pitch + grid item text)
├── layouts/
│   ├── index.html             frontpage: logo, magic command, feature grid
│   ├── _default/list.html     /blog/ index, newest-first, uses H1-title partial
│   ├── _default/single.html   post page, H1 title + .Date
│   └── partials/
│       ├── post-title.html    ← extract first "# " from .RawContent (the clever bit)
│       ├── header.html        logo + nav
│       └── footer.html
├── assets/css/main.css        dark, monospace, mascot-blue accents
└── static/                    jjay-persona.png, hero.png (copied from artwork/)

.github/workflows/pages.yml    Hugo build + Actions→Pages deploy, trigger push:main
```

## Risk ledger

| Risk | Mitigation |
|------|------------|
| Subpath `/jjay/` breaks CSS/images (works locally, 404 on Pages) | All refs via `relURL` / `.RelPermalink`; never a leading-slash literal. Verification task. |
| H1→title extraction fragility | Contained in one partial parsing `.RawContent`; **fallback** to prettified filename if no `# ` heading is found — degrades, never crashes. |
| A post with no leading `# ` | Same fallback path → filename-derived title. |
| Pages not set to "GitHub Actions" source | One-time repo setting; called out in tasks. |
| Org (`speclib`) restricts Pages/Actions | Verification task: confirm org allows Pages + Actions deploy. |
| Hugo version drift between local and CI | Pin Hugo version in the workflow. |

## Open questions (non-blocking)

- Exact feature-grid copy (one line per command) — authored during apply.
- RSS feed — Hugo gives it nearly free; revisit as a follow-up if wanted.
