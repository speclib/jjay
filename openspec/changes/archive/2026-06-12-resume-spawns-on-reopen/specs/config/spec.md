## ADDED Requirements

### Requirement: Agent commands resolve from a 3-layer, per-field config
jjay SHALL resolve each agent's `launch` and `resume` command templates from three sources, in priority order: the project file `<repo>/.jjay/config.yaml`, then the global file `~/.config/jjay/config.yaml`, then a Go built-in default. Resolution SHALL be **per field**: each of `launch` and `resume` is resolved independently as `project ?? global ?? builtin`, where a missing or empty string counts as unset. The resolver SHALL be defined over an ordered list of sources so that additional layers can be inserted later without changing call sites.

#### Scenario: Project overrides one field, inherits the rest
- **WHEN** `<repo>/.jjay/config.yaml` sets `agents.claude.launch` but not `agents.claude.resume`
- **THEN** the resolved `launch` is the project value
- **THEN** the resolved `resume` falls back to the global value if set, otherwise the built-in value

#### Scenario: No config files present
- **WHEN** neither `<repo>/.jjay/config.yaml` nor `~/.config/jjay/config.yaml` exists
- **THEN** both `launch` and `resume` resolve to the Go built-in defaults
- **THEN** spawning and reopening still work with zero configuration

#### Scenario: Present-but-blank field falls back
- **WHEN** a config file sets `agents.claude.resume` to an empty string
- **THEN** `resume` is treated as unset and falls back to the next source

### Requirement: Built-in defaults are the single source of truth, materialized on init
The Go built-in default profile SHALL be the canonical source for both the runtime fallback and the file `jjay init` seeds. `jjay init` SHALL write the built-in `agents` block into `<repo>/.jjay/config.yaml` so the values are visible and editable. Seeding SHALL be idempotent and non-destructive: an existing `.jjay/config.yaml` SHALL NOT be overwritten without `--force`. The built-in `claude` profile SHALL be `launch: claude "/opsx:apply {change}" --dangerously-skip-permissions --add-dir {wsdir}` and `resume: claude --resume --add-dir {wsdir}`.

#### Scenario: init seeds the project config
- **WHEN** `jjay init` runs in a repo with no `.jjay/config.yaml`
- **THEN** `<repo>/.jjay/config.yaml` is created containing `agents.claude.launch` and `agents.claude.resume`
- **THEN** the seeded values equal the Go built-in defaults

#### Scenario: init does not clobber an existing config
- **WHEN** `<repo>/.jjay/config.yaml` already exists and `jjay init` runs without `--force`
- **THEN** the existing file is left unchanged

### Requirement: jjay config is distinct from openspec config
The jjay config files (`<repo>/.jjay/config.yaml`, `~/.config/jjay/config.yaml`) SHALL be separate from `openspec/config.yaml`. jjay SHALL NOT write its `agents` block into `openspec/config.yaml`, and reading openspec's config SHALL NOT be required to resolve agent commands.

#### Scenario: openspec config untouched
- **WHEN** jjay resolves agent commands or `jjay init` seeds jjay config
- **THEN** `openspec/config.yaml` is neither read for nor written with jjay agent settings
