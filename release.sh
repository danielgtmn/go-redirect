#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Go Redirect Release Script${NC}"
echo "==========================="

CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${RED}Error: Must be on 'main' branch${NC}"
    exit 1
fi

if [ -n "$(git status --porcelain)" ]; then
    echo -e "${RED}Error: Working tree is not clean${NC}"
    git status --short
    exit 1
fi

LATEST_TAG=$(git tag --sort=-version:refname | head -n1 || echo "v0.0.0")
echo -e "${YELLOW}Latest tag: ${LATEST_TAG}${NC}"

echo -e "\n${YELLOW}Release type?${NC}"
echo "1) patch (1.0.0 → 1.0.1)"
echo "2) minor (1.0.0 → 1.1.0)"
echo "3) major (1.0.0 → 2.0.0)"
echo "4) custom"
read -p "Choose (1-4): " choice

case $choice in
    1)
        IFS='.' read -r major minor patch <<< "${LATEST_TAG#v}"
        NEW_VERSION="v${major}.${minor}.$((patch + 1))"
        ;;
    2)
        IFS='.' read -r major minor patch <<< "${LATEST_TAG#v}"
        NEW_VERSION="v${major}.$((minor + 1)).0"
        ;;
    3)
        IFS='.' read -r major minor patch <<< "${LATEST_TAG#v}"
        NEW_VERSION="v$((major + 1)).0.0"
        ;;
    4)
        read -p "Version (without 'v'): " custom_version
        NEW_VERSION="v${custom_version}"
        ;;
    *)
        echo -e "${RED}Invalid choice${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}New version: ${NEW_VERSION}${NC}"

if git tag --list | grep -q "^${NEW_VERSION}$"; then
    echo -e "${RED}Tag ${NEW_VERSION} already exists!${NC}"
    exit 1
fi

read -p "Create release ${NEW_VERSION}? (y/N): " confirm
if [[ ! $confirm =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Aborted${NC}"
    exit 0
fi

git tag -a "${NEW_VERSION}" -m "Release ${NEW_VERSION}"
git push origin "${NEW_VERSION}"

echo -e "\n${GREEN}Tag ${NEW_VERSION} pushed!${NC}"
echo -e "\nNext: Go to GitHub Actions → 'Create Release' → Run workflow"
