# Design: Setup Initial README

## Structure

```
README.md
├── Hero image (centered)
├── One-liner description
├── What jjay automates (manual workflow)
├── Prerequisites
├── Installation (placeholder)
├── CLI preview
├── Roadmap
├── Contributing
└── License
```

## Decisions

### Hero image rendering
Use HTML `<p align="center">` with `<img>` tag for centered display. Reference `artwork/hero.png` directly.

### Manual workflow section
Present the 6-step manual process as a numbered list with shell commands. This is the core of the README — it shows the pain and makes jjay's value obvious without a "why" section.

### CLI preview
Show planned commands in a fenced code block. Mark as preview/planned — nothing works yet.

### Roadmap
Bullet list, not a table or timeline. Features grouped loosely by priority. No dates — this is alpha.

### Tone
Technical, minimal, no fluff. No emojis in the README body. Let the hero image carry the personality.
