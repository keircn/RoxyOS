#!/bin/bash
set -e

cd "$(dirname "$0")"

REPO_NAME="roxyos"
REPO_DIR="repo"

echo "Building RoxyOS packages for repository..."
echo ""

mkdir -p "$REPO_DIR"

PACKAGES=(
  roxyos-assets
  roxyos-niri
  roxyos-waybar
  roxyos-hyprlock
  roxyos-kitty
  roxyos-ghostty
  roxyos-rofi
  roxyos-plymouth
  roxyos-fish
  roxyos-mako
  roxyos-ly
  roxyos-sddm
  roxyos-setup
  roxyos-hypr
  roxyos-fastfetch
  roxyos
)

for pkg in "${PACKAGES[@]}"; do
  echo "==> Building $pkg..."
  pkg_dir="packages/$pkg"

  if [[ ! -d "$pkg_dir" ]]; then
    echo "Warning: $pkg_dir not found, skipping"
    continue
  fi

  pushd "$pkg_dir" >/dev/null

  rm -rf pkg src *.pkg.tar.* 2>/dev/null || true

  makepkg -f

  mv *.pkg.tar.* "../../$REPO_DIR/" 2>/dev/null || true

  popd >/dev/null
  echo ""
done

echo "==> Creating repository database..."
cd "$REPO_DIR"
repo-add "${REPO_NAME}.db.tar.gz" *.pkg.tar.*

echo ""
echo "Repository created in $REPO_DIR/"
echo ""
echo "Files to upload:"
ls -la
