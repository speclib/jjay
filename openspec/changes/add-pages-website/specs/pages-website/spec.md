## ADDED Requirements

### Requirement: Devblog renders from devblog/posts directly
The site SHALL render the devblog from the `devblog/posts/` directory via a Hugo module mount, without copying or transforming the files. The `devblog/posts/*.md` files SHALL remain unmodified (no frontmatter added).

#### Scenario: Existing post appears on the site
- **WHEN** the site is built
- **THEN** each `devblog/posts/YYYY-MM-DD-slug.md` file appears as a post under `/blog/`
- **THEN** the source file is unchanged (no frontmatter, no edits)

#### Scenario: New post auto-publishes
- **WHEN** a new `devblog/posts/*.md` file is committed and pushed to `main`
- **THEN** the build includes it and the deployed `/blog/` lists it without any other change

### Requirement: Post title comes from the H1 heading
The site SHALL derive each post's title from the first `# ` heading in the post body, and its date from the `YYYY-MM-DD-` filename prefix. The title SHALL be consistent between the blog list page and the individual post page.

#### Scenario: Title and date resolved from a normal post
- **WHEN** a post named `2026-06-03-merge-command.md` begins with `# Merge, In One Word`
- **THEN** its title is "Merge, In One Word" on both the list and the post page
- **THEN** its date is 2026-06-03

#### Scenario: Post missing an H1 heading
- **WHEN** a post has no leading `# ` heading
- **THEN** the build does not fail
- **THEN** the title falls back to a prettified form of the filename slug

### Requirement: Compact frontpage for advanced developers
The frontpage SHALL present, compactly, the jjay logo, the "Control the flock." tagline, a short description, a single command that represents the spawn workflow, a feature grid covering the core commands, and links to the GitHub repository and the devblog.

#### Scenario: Frontpage content
- **WHEN** a visitor loads `/`
- **THEN** the `jjay-persona.png` logo and the "Control the flock." tagline are shown
- **THEN** a single `jjay spawn` command block is shown
- **THEN** a feature grid lists spawn, merge, cleanup, and session-open
- **THEN** links to the GitHub repository and to `/blog/` are present

### Requirement: Site builds for a project subpath
The site SHALL be configured with `baseURL` `https://speclib.github.io/jjay/` and SHALL reference all assets and internal links relative to that base so they resolve correctly under the `/jjay/` subpath.

#### Scenario: Assets resolve under the subpath
- **WHEN** the site is built with the project `baseURL` and served under `/jjay/`
- **THEN** CSS, images, and internal links resolve (no `/`-rooted paths that 404 under the subpath)

### Requirement: Automated build and deploy via GitHub Actions
A GitHub Actions workflow SHALL build the site with Hugo and deploy it to GitHub Pages using the Actions→Pages artifact mechanism, triggered on pushes to `main` that touch the site or the devblog posts.

#### Scenario: Push triggers deploy
- **WHEN** a commit touching `site/**` or `devblog/posts/**` is pushed to `main`
- **THEN** the workflow builds with the pinned Hugo version and deploys the result to GitHub Pages

#### Scenario: Unrelated push does not deploy
- **WHEN** a commit touches neither `site/**` nor `devblog/posts/**`
- **THEN** the Pages workflow does not run
