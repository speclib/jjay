## 1. Scaffold the Hugo site

- [ ] 1.1 Create `site/hugo.toml` with `baseURL = "https://speclib.github.io/jjay/"` and a module mount of `../devblog/posts` → `content/blog`
- [ ] 1.2 Copy `artwork/jjay-persona.png` and `artwork/hero.png` into `site/static/`
- [ ] 1.3 Create `site/content/_index.md` with frontpage copy (pitch, the `jjay spawn` command, feature-grid item text)

## 2. Layouts and styling

- [ ] 2.1 Create `site/layouts/partials/post-title.html` that extracts the first `# ` heading from `.RawContent`, with a fallback to a prettified filename slug when no heading exists
- [ ] 2.2 Create `site/layouts/index.html` (frontpage: logo, tagline, magic command, 4-item feature grid for spawn/merge/cleanup/session-open, GitHub + devblog links)
- [ ] 2.3 Create `site/layouts/_default/list.html` (devblog index, newest-first, using the post-title partial)
- [ ] 2.4 Create `site/layouts/_default/single.html` (post page: H1 title + date from filename)
- [ ] 2.5 Create `site/layouts/partials/header.html` and `footer.html` (logo + nav)
- [ ] 2.6 Create `site/assets/css/main.css` — dark, monospace, mascot-blue accent; reference all assets via `relURL`/`.RelPermalink` (no leading-slash literals)

## 3. GitHub Action

- [ ] 3.1 Create `.github/workflows/pages.yml`: pinned Hugo version, `actions/configure-pages` + `upload-pages-artifact` + `deploy-pages@v4`, permissions `pages: write` and `id-token: write`
- [ ] 3.2 Trigger on push to `main` with `paths: [site/**, devblog/posts/**]`

## 4. Repo / org configuration

- [ ] 4.1 Set repo Settings → Pages → Source = "GitHub Actions"
- [ ] 4.2 Confirm the `speclib` org allows Pages and Actions deploy

## 5. Verify

- [ ] 5.1 Build locally with the project `baseURL` and confirm CSS/images/links resolve under `/jjay/` (no 404s from `/`-rooted paths)
- [ ] 5.2 Confirm all 14 existing posts render under `/blog/`, newest first, with titles matching their `# ` headings
- [ ] 5.3 Confirm a post with no `# ` heading falls back to a filename title without failing the build
- [ ] 5.4 Push and confirm the Action deploys; load `https://speclib.github.io/jjay/` and `/jjay/blog/`
