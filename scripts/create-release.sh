#!/bin/bash

# Script to create a new release
# Usage: ./scripts/create-release.sh [version]

set -e

VERSION=${1:-"1.0.0"}
TAG="v${VERSION}"

echo "🚀 Creating release ${TAG}..."

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo "❌ Not in a git repository!"
    exit 1
fi

# Check if tag already exists
if git tag -l | grep -q "^${TAG}$"; then
    echo "❌ Tag ${TAG} already exists!"
    exit 1
fi

# Build packages
echo "📦 Building packages..."
./scripts/build-packages.sh

# Create and push tag
echo "🏷️  Creating tag ${TAG}..."
git tag -a "${TAG}" -m "Release ${TAG}"

echo "📤 Pushing tag to remote..."
git push origin "${TAG}"

echo "✅ Release ${TAG} created successfully!"
echo ""
echo "📋 Next steps:"
echo "1. GitHub Actions will automatically build and upload packages"
echo "2. Check the Actions tab in your GitHub repository"
echo "3. Once complete, packages will be available in the Releases section"
echo ""
echo "🎉 Release process initiated!"
