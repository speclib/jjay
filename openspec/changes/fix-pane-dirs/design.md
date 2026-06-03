## Context

After spawn, the right pane (shell) ends up in the wrong directory. Root cause: `send-keys "cd <wsDir>"` races with shell initialization. The fix was designed in the spawn-config change but was dropped during a merge conflict.

## Goals / Non-Goals

**Goals:**
- Fix pane working directories using tmux's `-c` flag
- Add integration test that catches this regression
- Add Makefile target for running integration tests

**Non-Goals:**
- Changing any other spawn behavior
- Unit tests for tmux commands (integration test covers this)

## Decisions

### Use tmux `-c` flag for working directories

Replace `send-keys cd` with `-c` on window/pane creation:

```
Before:
  tmux new-window -d -n ws-X
  tmux send-keys -t ws-X "cd /path && agent-cmd" Enter
  tmux split-window -h -t ws-X
  tmux send-keys -t ws-X.1 "cd /path" Enter

After:
  tmux new-window -d -n ws-X -c /path
  tmux send-keys -t ws-X "agent-cmd" Enter
  tmux split-window -h -t ws-X -c /path
  (no send-keys for right pane — it starts in /path)
```

The agent command no longer needs `cd <wsDir> &&` prefix since the window already starts in the workspace dir.

### Integration test structure

```
Test setup:
  1. Create temp dir for test workspace root
  2. Init jj repo in temp dir
  3. Create openspec change dir (mkdir, not full openspec init)
  4. Create dedicated tmux session: jjay-test-<random>
  5. Build fake-agent.sh path

Test body:
  1. Spawn(changeName, opts) with fake agent, test session, temp root
  2. Sleep briefly for agent to run
  3. Assert: tmux window exists
  4. Assert: both panes report correct working dir
  5. Assert: jj workspace exists
  6. Assert: workspace dir exists with agent marker file
  7. Cleanup(changeName, opts) with test session, temp root
  8. Assert: tmux window gone
  9. Assert: jj workspace gone
  10. Assert: workspace dir gone

Teardown (defer):
  - Kill test tmux session
  - Remove temp dirs
```

### Verify pane directory via tmux

```bash
tmux display-message -p -t "session:window.0" '#{pane_current_path}'
```

Returns the absolute path of the pane's current working directory. Compare against the expected workspace dir.

### Fake agent: testdata/fake-agent.sh

```bash
#!/bin/sh
echo "fake agent: $@" > agent-was-here.txt
```

Creates a marker file so the test can verify the agent ran in the right directory.

## Risks / Trade-offs

- [Integration test timing] → Agent runs async in tmux pane. Need brief sleep (1-2s) for marker file. Poll with timeout instead of fixed sleep.
- [tmux/jj required] → Skip test if not available (`t.Skip`).
