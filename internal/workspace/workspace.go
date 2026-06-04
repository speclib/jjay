package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WindowName returns the tmux window name for a given change.
func WindowName(changeName string) string {
	return "ws-" + changeName
}

// MainRepoRoot returns the root directory of the main jj working copy (the
// `default` workspace), resolved from `start` (or any of its ancestors).
//
// jj stores a `.jj/repo` entry in every workspace: in the main working copy it
// is a real directory, while in a child workspace it is a file whose contents
// point (relative to `.jj/`) at the main copy's `.jj/repo`. We walk up from
// `start` to the nearest `.jj`, then follow the pointer if present, so the
// result is the same whether jjay is run from the main repo or any child
// workspace.
func MainRepoRoot(start string) (string, error) {
	jjDir, err := findJJDir(start)
	if err != nil {
		return "", err
	}

	repoEntry := filepath.Join(jjDir, "repo")
	info, err := os.Stat(repoEntry)
	if err != nil {
		return "", fmt.Errorf("failed to stat %s: %w", repoEntry, err)
	}

	// Main working copy: .jj/repo is a directory; its parent's parent is root.
	if info.IsDir() {
		return filepath.Dir(jjDir), nil
	}

	// Child workspace: .jj/repo is a file pointing at <main>/.jj/repo.
	data, err := os.ReadFile(repoEntry)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", repoEntry, err)
	}
	pointer := strings.TrimSpace(string(data))
	mainRepoEntry := pointer
	if !filepath.IsAbs(mainRepoEntry) {
		mainRepoEntry = filepath.Join(jjDir, pointer)
	}
	// mainRepoEntry == <main>/.jj/repo → strip /repo and /.jj to get <main>.
	mainRoot := filepath.Dir(filepath.Dir(mainRepoEntry))
	abs, err := filepath.Abs(mainRoot)
	if err != nil {
		return "", fmt.Errorf("failed to resolve main repo root: %w", err)
	}
	return filepath.Clean(abs), nil
}

// findJJDir walks up from start looking for a `.jj` directory, returning its
// absolute path.
func findJJDir(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, ".jj")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no .jj directory found in %s or its ancestors", start)
		}
		dir = parent
	}
}

// RelativeToMain returns target expressed relative to the main repo root. If a
// relative path cannot be computed (e.g. different volumes), the absolute path
// is returned unchanged.
func RelativeToMain(mainRoot, target string) string {
	rel, err := filepath.Rel(mainRoot, target)
	if err != nil {
		return target
	}
	return rel
}

// WorkspaceDir returns the absolute path for the workspace directory.
// If root is empty, uses the default: ../<project-name>-workspaces/<change-name>.
// If root is set, uses: <root>/<change-name>.
func WorkspaceDir(changeName, root string) (string, error) {
	if root != "" {
		absPath, err := filepath.Abs(filepath.Join(root, changeName))
		if err != nil {
			return "", fmt.Errorf("failed to resolve workspace path: %w", err)
		}
		return absPath, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	return workspaceDirUnder(cwd, changeName)
}

// WorkspaceDirFrom resolves a workspace directory anchored on the main repo
// root rather than the current directory, so it is correct even when jjay is
// invoked from inside a child workspace. If root is set it takes precedence
// (same as WorkspaceDir); otherwise the default convention is applied relative
// to mainRoot: <mainRoot>/../<mainRoot-name>-workspaces/<change>.
func WorkspaceDirFrom(mainRoot, changeName, root string) (string, error) {
	if root != "" {
		return WorkspaceDir(changeName, root)
	}
	return workspaceDirUnder(mainRoot, changeName)
}

// workspaceDirUnder applies the default workspace-dir convention relative to
// base: ../<base-name>-workspaces/<change>.
func workspaceDirUnder(base, changeName string) (string, error) {
	projectName := filepath.Base(base)
	absPath, err := filepath.Abs(filepath.Join(base, "..", projectName+"-workspaces", changeName))
	if err != nil {
		return "", fmt.Errorf("failed to resolve workspace path: %w", err)
	}
	return absPath, nil
}
