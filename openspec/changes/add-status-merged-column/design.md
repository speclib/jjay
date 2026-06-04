## Context

`internal/status.List()` builds `[]Spawn{Change, WSDir, Attached, Archived, Tasks}` by joining `jj workspace list` with `tmux list-windows`, then `Render` prints a tabwriter table headed `CHANGE WORKSPACE TASKS ARCHIVED STATUS`. This change adds a `Merged` field + `MERGED` column and renames the tmux-state header `STATUS → TMUX`. All derivation stays live (ADR-006); no state file.

A spike during exploration found the merged signal is subtler than it looks: because `rebase-before-merge` rebases every spawn onto `main`, naive ancestor checks (`<change>@ & ::main`, or filtering empty commits) give false positives — they reported active, unmerged spawns as merged. The signal that matched reality (archived changes = merged, active changes = not) was **"does the workspace have commits ahead of main"**: `main..<change>@` empty ⇒ merged.

## Goals / Non-Goals

**Goals:**
- Add a `MERGED` (yes/no) column, derived live from jj.
- Rename the tmux-state column `STATUS → TMUX`.
- Keep `ARCHIVED`; final order `CHANGE WORKSPACE TASKS TMUX MERGED ARCHIVED`.

- Add `jjay status` to the lifecycle integration test (folds in jjay-y5yn), giving the MERGED edge cases an end-to-end check against a real jj repo.

**Non-Goals (→ jjay-atq6):**
- `AGENT` column, `STATUS` busy/finished, `--columns`/view configurability.
- Any persisted or agent-emitted state.
- Auto-cleanup of merged spawns (→ jjay-0yko); this only surfaces the state.

## Decisions

- **`Merged bool` on `Spawn`**, computed in `List`/`join` (one jj call per spawn, or one batched revset). Render adds the `MERGED` column and renames the header.
- **Leading detection candidate:** `main..<change>@` is empty ⇒ merged. The implementer SHOULD confirm edge cases before locking it in:
  - **Empty workspace** (no work yet): `main..<change>@` empty ⇒ would read as "merged". Is that right, or should an empty/never-worked spawn show `MERGED=no`? (Likely treat empty-and-no-work as `no` — it isn't "done".) The existing `checkWorkspaceEmpty` logic in `merge.go` is a reference.
  - **Just-rebased, not yet merged:** after `merge`'s pre-rebase step a workspace sits on `main` but still has its own commits ahead — `main..<change>@` should be non-empty ⇒ `no`. Verify.
  - **Merge commit semantics:** `merge` does `jj new main <change>@`; confirm the revset reflects post-merge state correctly for a workspace that lingers.
- **No new tool/package** — jj only, inside `internal/status`.
- **Tolerance:** if the jj revset can't be evaluated, default `MERGED=no` (don't fail status), consistent with status's tmux-tolerance stance.
- **Integration test (jjay-y5yn):** add a `status` subtest to `TestFullLifecycle` between `spawn` and `cleanup`, reusing the existing `testEnv`. Call `status.List(env.SessionName, env.WsRoot)` + `Render` (or the binary) and assert the live spawn shows up with `MERGED=no` and the TMUX/MERGED columns. This is the natural oracle for the merged edge cases — a freshly spawned workspace is unmerged by construction. Deeper merged=yes coverage stays in the `internal/merge` integration suite, which already drives real merges.

## Risks / Trade-offs

- **Revset correctness is the whole risk.** The spike showed easy ways to get false positives; design deliberately leaves the implementer room to validate the candidate against archived-vs-active spawns (the spike's oracle) and the empty-workspace edge before finalizing.
- **Per-spawn jj call** adds a little latency to `status` (already does per-spawn tasks.md reads, so consistent). Could batch into one revset over all workspace heads if needed.
- **MERGED vs ARCHIVED confusion** for users — mitigated by keeping both columns and the proposal's framing ("merged but not archived = ready to cleanup").
