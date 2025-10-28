#!/bin/bash

# MovieStream Desktop Release Script
# This script helps you create a new release automatically

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}MovieStream Desktop - Release Creator${NC}"
echo "=========================================="
echo ""

# Check if git is clean
if [[ -n $(git status -s) ]]; then
    echo -e "${RED}Error: Working directory is not clean!${NC}"
    echo "Please commit or stash your changes first."
    exit 1
fi

# Get current branch
BRANCH=$(git branch --show-current)
echo -e "Current branch: ${YELLOW}$BRANCH${NC}"

# Get the last tag
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo -e "Last release: ${YELLOW}$LAST_TAG${NC}"
echo ""

# Ask for version
echo "Enter new version (format: v1.0.0):"
read -r VERSION

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}Error: Invalid version format!${NC}"
    echo "Version must be in format: v1.0.0"
    exit 1
fi

# Check if tag already exists
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo -e "${RED}Error: Tag $VERSION already exists!${NC}"
    exit 1
fi

echo ""
echo "Creating release ${GREEN}$VERSION${NC}"
echo "This will:"
echo "  1. Create and push a git tag"
echo "  2. Trigger GitHub Actions to:"
echo "     - Create a GitHub release"
echo "     - Build Windows .exe"
echo "     - Build Linux binary"
echo "     - Build macOS binaries (Intel + ARM)"
echo "     - Attach all binaries to the release"
echo ""
echo -e "${YELLOW}Continue? (y/n)${NC}"
read -r CONFIRM

if [[ $CONFIRM != "y" && $CONFIRM != "Y" ]]; then
    echo "Cancelled."
    exit 0
fi

# Create and push tag
echo ""
echo "Creating tag..."
git tag -a "$VERSION" -m "Release $VERSION"

echo "Pushing tag to remote..."
git push origin "$VERSION"

echo ""
echo -e "${GREEN}✓ Success!${NC}"
echo ""
echo "Release tag created and pushed!"
echo ""
echo "GitHub Actions is now:"
echo "  → Creating the release"
echo "  → Building binaries for all platforms"
echo "  → Uploading binaries to the release"
echo ""
echo "View progress at:"
echo -e "${YELLOW}https://github.com/zerdnem/moviestream-desktop/actions${NC}"
echo ""
echo "Release will be available at:"
echo -e "${YELLOW}https://github.com/zerdnem/moviestream-desktop/releases/tag/$VERSION${NC}"
echo ""
echo "This usually takes 5-10 minutes to complete."

