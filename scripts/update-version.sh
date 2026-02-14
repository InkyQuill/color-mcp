#!/bin/bash
set -e

VERSION=$1
if [ -z "$VERSION" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

# Update version in main.go
sed -i "s/serverVersion = \".*\"/serverVersion = \"${VERSION}\"/" main.go

echo "Updated main.go to version ${VERSION}"
