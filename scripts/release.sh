#!/usr/bin/env bash
set -euo pipefail

# jjay release script — jj only, no git commands

# --- Safety checks ---

if ! command -v gum &>/dev/null; then
  echo "Error: gum is required. Install it or enter the nix devShell."
  exit 1
fi

if ! command -v jj &>/dev/null; then
  echo "Error: jj is required."
  exit 1
fi

# Check working copy is clean (empty)
if jj status 2>/dev/null | grep -q "Working copy changes:"; then
  echo "Error: working copy has uncommitted changes. Run 'jj new' first."
  exit 1
fi

# Check main bookmark points to current parent
MAIN_COMMIT=$(jj log -r "main" --no-graph -T 'commit_id.short(12)' 2>/dev/null)
PARENT_COMMIT=$(jj log -r "@-" --no-graph -T 'commit_id.short(12)' 2>/dev/null)
if [[ "$MAIN_COMMIT" != "$PARENT_COMMIT" ]]; then
  echo "Error: not on main. Current parent is not the main bookmark."
  echo "  main:   $MAIN_COMMIT"
  echo "  parent: $PARENT_COMMIT"
  exit 1
fi

if ! grep -q "## Unreleased" CHANGELOG.md; then
  echo "Error: CHANGELOG.md has no '## Unreleased' section."
  exit 1
fi

# --- Read current version ---

CURRENT=$(cat VERSION | tr -d '[:space:]')
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT"

echo "Current version: $CURRENT"

# --- Interactive version bump ---

BUMP=$(gum choose "patch" "minor" "major")

case "$BUMP" in
  major) MAJOR=$((MAJOR + 1)); MINOR=0; PATCH=0 ;;
  minor) MINOR=$((MINOR + 1)); PATCH=0 ;;
  patch) PATCH=$((PATCH + 1)) ;;
esac

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
TAG="v${NEW_VERSION}"

# Check tag doesn't exist (via git tags, which jj exports to)
if git tag -l "$TAG" 2>/dev/null | grep -q "$TAG"; then
  echo "Error: tag $TAG already exists."
  exit 1
fi

echo "Bumping to: $NEW_VERSION (tag: $TAG)"
gum confirm "Proceed?" || exit 0

# --- Update VERSION file ---

echo "$NEW_VERSION" > VERSION

# --- Update CHANGELOG.md ---

DATE=$(date +%Y-%m-%d)
sed -i "s/## Unreleased/## Unreleased\n\n## ${NEW_VERSION} - ${DATE}/" CHANGELOG.md

# --- Update nix vendorHash ---

if command -v nix &>/dev/null; then
  echo "Updating nix vendorHash..."

  FAKE_HASH="sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
  sed -i "s|vendorHash = \".*\"|vendorHash = \"${FAKE_HASH}\"|" flake.nix

  CORRECT_HASH=$(nix build 2>&1 | grep -oP 'got:\s+\Ksha256-[A-Za-z0-9+/=]+' | head -1) || true

  if [[ -n "$CORRECT_HASH" ]]; then
    sed -i "s|vendorHash = \".*\"|vendorHash = \"${CORRECT_HASH}\"|" flake.nix
    echo "Updated vendorHash to: $CORRECT_HASH"
  else
    echo "Warning: could not determine correct vendorHash."
    jj restore --from @- -- flake.nix
  fi
else
  echo "Warning: nix not installed, skipping vendorHash update."
fi

# --- Describe, bookmark, tag, push ---

jj describe -m "release: v${NEW_VERSION}"
jj bookmark set main -r @
jj new

# Create git tag (jj doesn't have native tags, use git via jj's colocated repo)
git tag -a "$TAG" -m "Release ${NEW_VERSION}"

echo "Pushing..."
jj git push
git push origin "$TAG"

echo "Done! Release $TAG has been pushed."
echo "GitHub Actions will create the release automatically."
