#!/usr/bin/env bash
set -euo pipefail

# --- Safety checks ---

if ! command -v gum &>/dev/null; then
  echo "Error: gum is required. Install it or enter the nix devShell."
  exit 1
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Error: working tree is dirty. Commit or stash changes first."
  exit 1
fi

BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ "$BRANCH" != "main" ]]; then
  echo "Error: not on main branch (currently on $BRANCH)."
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

if git tag -l "$TAG" | grep -q "$TAG"; then
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

  # Set a fake hash to trigger the build error with the correct hash
  FAKE_HASH="sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
  sed -i "s|vendorHash = \".*\"|vendorHash = \"${FAKE_HASH}\"|" flake.nix

  # Build and capture the expected hash from stderr
  CORRECT_HASH=$(nix build 2>&1 | grep -oP 'got:\s+\Ksha256-[A-Za-z0-9+/=]+' | head -1) || true

  if [[ -n "$CORRECT_HASH" ]]; then
    sed -i "s|vendorHash = \".*\"|vendorHash = \"${CORRECT_HASH}\"|" flake.nix
    echo "Updated vendorHash to: $CORRECT_HASH"
  else
    echo "Warning: could not determine correct vendorHash. Restoring original."
    git checkout -- flake.nix
  fi
else
  echo "Warning: nix not installed, skipping vendorHash update."
fi

# --- Git commit, tag, push ---

git add VERSION CHANGELOG.md flake.nix
git commit -m "release: v${NEW_VERSION}"
git tag -a "$TAG" -m "Release ${NEW_VERSION}"

echo "Pushing commit and tag..."
git push origin main
git push origin "$TAG"

echo "Done! Release $TAG has been pushed."
echo "GitHub Actions will create the release automatically."
