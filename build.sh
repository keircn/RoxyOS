#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "RoxyOS Package Builder"
echo "======================"
echo ""

echo "Select display manager:"
echo "  1) ly (console-based, lightweight)"
echo "  2) sddm (graphical, feature-rich)"
echo "  3) none (skip display manager)"
echo ""
read -rp "Choice [1-3]: " dm_choice

case "$dm_choice" in
1) DM_PACKAGE="roxyos-ly" ;;
2) DM_PACKAGE="roxyos-sddm" ;;
3) DM_PACKAGE="" ;;
*)
  echo "Invalid choice"
  exit 1
  ;;
esac

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
  roxyos-setup
)

if [[ -n "$DM_PACKAGE" ]]; then
  PACKAGES+=("$DM_PACKAGE")
fi

PACKAGES+=(roxyos)

echo ""
echo "Building RoxyOS packages..."
echo "Packages to install: ${PACKAGES[*]}"
echo ""

for pkg in "${PACKAGES[@]}"; do
  echo "==> Building $pkg..."
  cd "packages/$pkg"
  makepkg -si --noconfirm
  cd ../..
  echo ""
done

echo "All packages built and installed successfully!"

if [[ -n "$DM_PACKAGE" ]]; then
  if [[ "$DM_PACKAGE" == "roxyos-ly" ]]; then
    echo "To enable ly: sudo systemctl enable ly.service"
  else
    echo "To enable sddm: sudo systemctl enable sddm.service"
  fi
fi
