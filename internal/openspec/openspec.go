// Package openspec is a lean reader over the openspec CLI. It runs
// `openspec list --json` and returns the change names, nothing more — no tmux,
// no task files. It is the single source for "what openspec changes exist",
// shared by spawn's precondition check and completion's candidate logic
// (see ADR-009).
package openspec

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type change struct {
	Name string `json:"name"`
}

type list struct {
	Changes []change `json:"changes"`
}

// ChangeNames returns the names of all openspec changes by parsing
// `openspec list --json`. It returns an error if openspec cannot be run or its
// output cannot be parsed; callers that must not fail (e.g. shell completion)
// should treat any error as "no candidates".
func ChangeNames() ([]string, error) {
	out, err := exec.Command("openspec", "list", "--json").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list openspec changes: %w", err)
	}
	return parseChangeNames(out)
}

// parseChangeNames extracts change names from `openspec list --json` output.
// Split out from ChangeNames so the parse is unit-testable without invoking the
// openspec binary.
func parseChangeNames(out []byte) ([]string, error) {
	var parsed list
	if err := json.Unmarshal(out, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse openspec output: %w", err)
	}

	names := make([]string, 0, len(parsed.Changes))
	for _, c := range parsed.Changes {
		names = append(names, c.Name)
	}
	return names, nil
}
