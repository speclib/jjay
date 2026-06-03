---
# jjay-v4oe
title: 'release script: add progress feedback during nix vendorHash update'
status: todo
type: task
priority: normal
created_at: 2026-06-03T15:53:23Z
updated_at: 2026-06-03T19:06:31Z
parent: jjay-qltp
---

The nix vendorHash update step in scripts/release.sh runs nix build silently while piping to grep. This can take minutes with no feedback. Add status messages like 'Building with fake hash to determine correct vendorHash (this may take a minute)...' and show nix build progress or a spinner.
