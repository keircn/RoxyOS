#!/bin/bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
    echo "Usage: $0 <version> [pkgrel]"
    echo "Example: $0 1.1.0"
    echo "Example: $0 1.1.0 2"
    exit 1
fi

VERSION="$1"
PKGREL="${2:-1}"

if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must be in format X.Y.Z (e.g., 1.0.0)"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PACKAGES_DIR="$SCRIPT_DIR/packages"

count=0
for pkgbuild in "$PACKAGES_DIR"/*/PKGBUILD; do
    if [[ -f "$pkgbuild" ]]; then
        sed -i "s/^pkgver=.*/pkgver=$VERSION/" "$pkgbuild"
        sed -i "s/^pkgrel=.*/pkgrel=$PKGREL/" "$pkgbuild"
        pkg=$(basename "$(dirname "$pkgbuild")")
        echo "Updated $pkg to $VERSION-$PKGREL"
        count=$((count + 1))
    fi
done

echo ""
echo "Bumped $count packages to version $VERSION-$PKGREL"
