# RoxyOS

Roxy themed metapackage for Arch Linux with Niri configs

## Installation

```bash
# Install RoxyOS metapackage
paru -S roxyos

# Run setup script
sudo roxyos-setup

# Follow the prompts to complete setup
```

## Post-Installation

After running `roxyos-setup`:

1. Log out and select "Niri" from your display manager
2. Or run `niri` from tty

## Keybindings

| Shortcut | Action |
|----------|--------|
| `Super + T` | Open terminal (kitty) |
| `Super + D` | Open application launcher (rofi) |
| `Super + B` | Open browser |
| `Super + L` | Lock screen |
| `Super + V` | Clipboard history |
| `Super + Q` | Close window |
| `Super + F` | Toggle fullscreen |
| `Super + Space` | Toggle floating |
| `Super + 1-0` | Switch to workspace 1-10 |
| `Super + Shift + 1-10` | Move window to workspace |
| `Super + Arrow keys` | Focus window |
| `Super + H/J/K/L` | Focus window (vim keys) |
| `Super + Ctrl + Arrow keys` | Resize window |
| `Super + S` | Toggle special workspace |

## Updating Configurations

```bash
# Apply all configurations
sudo roxyos-apply all

# Apply specific configuration
sudo roxyos-apply niri
sudo roxyos-apply waybar
sudo roxyos-apply kitty
sudo roxyos-apply rofi
sudo roxyos-apply fish
```

## Customization

### Theming

Edit `/etc/roxyos/setup.conf` to customize:

- Color scheme
- Display manager choice
- Plymouth settings
- Keybindings

### Wallpapers

Place your own wallpapers in `~/.local/share/wallpapers/roxyos/`

## Requirements

- Arch Linux or Arch-based distribution
- Paru (or yay if you simply must be weird)

## Subpackages

- `roxyos-niri` - Niri configuration
- `roxyos-waybar` - Waybar configuration
- `roxyos-hyprlock` - Hyprlock configuration
- `roxyos-kitty` - Kitty terminal configuration
- `roxyos-rofi` - Rofi launcher configuration
- `roxyos-plymouth` - Plymouth boot theme
- `roxyos-assets` - Visual assets (wallpapers, themes)
- `roxyos-fish` - Fish shell configuration
- `roxyos-mako` - Mako notification configuration

## License

RNC-1 License ([LICENSE](LICENSE))

## Support

Issues and suggestions welcome at the GitHub repository.
