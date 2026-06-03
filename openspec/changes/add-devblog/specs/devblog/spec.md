## ADDED Requirements

### Requirement: Blog artifact in schema
The `spec-driven-with-adr` schema SHALL include a `blog` artifact that generates files at `blog/<slug>.md` inside the change directory. At archive time, blog posts SHALL be synced to `devblog/posts/<date>-<slug>.md`.

#### Scenario: Schema includes blog artifact
- **WHEN** a change is created using the `spec-driven-with-adr` schema
- **THEN** the artifact graph includes `blog` alongside `proposal`, `specs`, `design`, `adr`, and `tasks`

#### Scenario: Blog artifact depends on proposal
- **WHEN** `openspec status` is checked for a new change
- **THEN** the `blog` artifact requires `proposal`

### Requirement: Blog post persona
All blog posts SHALL be narrated by Kaa, the project mascot (a Eurasian Jay / vlaamse gaai in a jujutsu gi). Posts SHALL use first person ("I", "my flock"), a brutaal and bossy tone, and be transparently AI-generated.

#### Scenario: Post voice
- **WHEN** a blog post is generated
- **THEN** it uses first person from Kaa's perspective
- **THEN** the tone is confident, direct, and playful

### Requirement: Blog posts are standalone
Each blog post SHALL cover one change and stand on its own. Posts SHALL be short and punchy, focusing on what new capabilities were added and why.

#### Scenario: Post covers one change
- **WHEN** a blog post is generated for change `spawn-command`
- **THEN** it covers only the spawn command's capabilities and motivation
- **THEN** a reader can understand it without reading other posts

### Requirement: Blog post location
Blog drafts SHALL live in the change directory at `blog/<slug>.md` during development. At archive time, they SHALL be synced to `devblog/posts/<date>-<slug>.md`. Posts SHALL NOT be deleted from the archive.

#### Scenario: Draft location during change
- **WHEN** a blog post is created for change `spawn-command`
- **THEN** the draft exists at `openspec/changes/spawn-command/blog/spawn-command.md`

#### Scenario: Post synced at archive
- **WHEN** the change is archived
- **THEN** the post is copied to `devblog/posts/<date>-spawn-command.md`

### Requirement: Devblog README
A `devblog/README.md` SHALL exist introducing Kaa and explaining that the blog is AI-generated raw material.

#### Scenario: README exists
- **WHEN** the devblog directory is inspected
- **THEN** `devblog/README.md` introduces Kaa and the blog concept
