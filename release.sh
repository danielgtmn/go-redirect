#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Caddy Redirect Release Script${NC}"
echo "=================================="

# Check if we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${RED}❌ Error: Must be on 'main' branch to create releases${NC}"
    echo -e "Current branch: ${CURRENT_BRANCH}"
    echo -e "Run: git checkout main"
    exit 1
fi

# Check if working tree is clean
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${RED}❌ Error: Working tree is not clean${NC}"
    echo -e "Please commit or stash your changes first"
    git status --short
    exit 1
fi

# Get latest version tag
LATEST_TAG=$(git tag --sort=-version:refname | head -n1 || echo "v0.0.0")
echo -e "${YELLOW}Latest tag: ${LATEST_TAG}${NC}"

# Ask for new version
echo -e "\n${YELLOW}What type of release?${NC}"
echo "1) patch (bug fixes) - 1.0.0 → 1.0.1"
echo "2) minor (new features) - 1.0.0 → 1.1.0"
echo "3) major (breaking changes) - 1.0.0 → 2.0.0"
echo "4) custom version"
read -p "Choose (1-4): " choice

case $choice in
    1)
        # Extract version numbers
        IFS='.' read -r major minor patch <<< "${LATEST_TAG#v}"
        new_patch=$((patch + 1))
        NEW_VERSION="v${major}.${minor}.${new_patch}"
        ;;
    2)
        IFS='.' read -r major minor patch <<< "${LATEST_TAG#v}"
        new_minor=$((minor + 1))
        NEW_VERSION="v${major}.${new_minor}.0"
        ;;
    3)
        IFS='.' read -r major minor patch <<< "${LATEST_TAG#v}"
        new_major=$((major + 1))
        NEW_VERSION="v${new_major}.0.0"
        ;;
    4)
        read -p "Enter custom version (without 'v'): " custom_version
        NEW_VERSION="v${custom_version}"
        ;;
    *)
        echo -e "${RED}❌ Invalid choice${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}New version: ${NEW_VERSION}${NC}"

# Check if tag already exists
if git tag --list | grep -q "^${NEW_VERSION}$"; then
    echo -e "${RED}❌ Error: Tag ${NEW_VERSION} already exists!${NC}"
    echo -e "Choose a different version or delete the existing tag:"
    echo -e "git tag -d ${NEW_VERSION} && git push origin :refs/tags/${NEW_VERSION}"
    exit 1
fi

# Confirm
read -p "Create release ${NEW_VERSION}? (y/N): " confirm
if [[ ! $confirm =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Aborted${NC}"
    exit 0
fi

# Create annotated tag
echo -e "\n${YELLOW}Creating git tag...${NC}"
git tag -a "${NEW_VERSION}" -m "Release ${NEW_VERSION}"

# Push tag to trigger release
echo -e "\n${YELLOW}Pushing tag to GitHub...${NC}"
git push origin "${NEW_VERSION}"

echo -e "\n${GREEN}✅ Tag ${NEW_VERSION} created and pushed!${NC}"
echo -e "\n${BLUE}What happens next:${NC}"
echo "1. GitHub Actions will build and push Docker images automatically"
echo "2. Go to GitHub Actions → 'Create Release' workflow"
echo "3. Click 'Run workflow' and enter: ${NEW_VERSION#v}"
echo "4. Docker images will be available at:"
echo "   - ghcr.io/danielgtmn/caddy-redirect:${NEW_VERSION#v}"
echo "   - ghcr.io/danielgtmn/caddy-redirect:latest"
echo -e "\n${YELLOW}Monitor the progress at:${NC}"
echo "https://github.com/danielgtmn/caddy-redirect/actions"
