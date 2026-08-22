#!/bin/bash
# Automates the add-on release ritual: tag the current commit, verify the
# tag actually resolves as a Go module version, bump addon/spanet/config.yaml
# and Dockerfile to match, verify the add-on image still builds, commit the
# bump, and publish a GitHub release (using the matching CHANGELOG.md
# section as its notes) - which also triggers release-binaries.yaml since
# that fires on `release: published`. Leaves the final `git push` to you.
#
# Preconditions this enforces:
#   - the gh CLI is installed and authenticated (gh auth login)
#   - working tree is clean (the fix itself is already committed)
#   - the tag doesn't already exist
#   - addon/spanet/CHANGELOG.md already has a "## <version>" section
#
# Usage: scripts/release-addon.sh 1.0.3
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

VERSION="${1:-}"
if [ -z "$VERSION" ]; then
	echo "usage: $0 <version, e.g. 1.0.3>" >&2
	exit 1
fi

if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
	echo "version must look like X.Y.Z (got: $VERSION)" >&2
	exit 1
fi

TAG="v$VERSION"

if ! command -v gh >/dev/null 2>&1; then
	echo "gh CLI is not installed (brew install gh)" >&2
	exit 1
fi

if ! gh auth status >/dev/null 2>&1; then
	echo "gh is not authenticated - run 'gh auth login' first" >&2
	exit 1
fi

if [ -n "$(git status --porcelain)" ]; then
	echo "working tree is not clean - commit the fix first" >&2
	exit 1
fi

if git rev-parse "$TAG" >/dev/null 2>&1; then
	echo "tag $TAG already exists" >&2
	exit 1
fi

if ! grep -qF "## $VERSION" addon/spanet/CHANGELOG.md; then
	echo "addon/spanet/CHANGELOG.md has no '## $VERSION' section - add one, commit it with the fix, then re-run" >&2
	exit 1
fi

# Everything between "## $VERSION" and the next "## " heading (or EOF).
NOTES="$(awk -v ver="## $VERSION" '
	$0 == ver { found=1; next }
	found && /^## / { exit }
	found { print }
' addon/spanet/CHANGELOG.md)"

echo "==> tagging $TAG at $(git rev-parse --short HEAD)"
git tag -a "$TAG" -m "$(git log -1 --pretty=%s)"
git push origin "$TAG"

echo "==> confirming the module resolves at $TAG"
GOFLAGS=-mod=mod go list -m "github.com/freman/spanet@$TAG"

echo "==> bumping addon/spanet/config.yaml and Dockerfile to $VERSION"
sed -i '' "s/^version: \".*\"/version: \"$VERSION\"/" addon/spanet/config.yaml
sed -i '' "s#cmd/spalink@v[0-9.]*#cmd/spalink@$TAG#" addon/spanet/Dockerfile

echo "==> verifying the add-on image still builds"
docker build --no-cache -t "spanet-addon:$TAG" addon/spanet >/dev/null
docker rmi "spanet-addon:$TAG" >/dev/null

echo "==> committing the bump"
git add addon/spanet/config.yaml addon/spanet/Dockerfile
git commit -m "Bump add-on to $TAG"

echo "==> publishing the GitHub release"
gh release create "$TAG" --title "$TAG" --notes "$NOTES"

echo
echo "Done. Review with 'git show', then push:"
echo "  git push"
