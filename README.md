# RoxyOS

<p align="center">
  <img src="" alt="RoxyOS Banner">
</p>

<p align="center">
  <strong>Roxy on Arch Linux</strong>
</p>

---

## Features

- Niri Window Manager configuration
- Waybar Status Bar
- Hyprlock screen locker
- Kitty terminal configuration
- Rofi application launcher
- Plymouth boot splash
- Fish shell with Starship prompt
- Mako notifications

## Installation

```bash
# Clone and build
git clone https://github.com/keircn/RoxyOS.git
cd RoxyOS
makepkg -si

# Run setup
sudo roxyos-setup
```

## Keybindings

| Shortcut | Action |
|----------|--------|
| Super + T | Terminal |
| Super + D | Launcher |
| Super + L | Lock screen |
| Super + Q | Close window |
| Super + F | Fullscreen |
| Super + Space | Float |
| Super + 1-0 | Workspace |

## Packages

- roxyos (main metapackage)
- roxyos-niri, roxyos-waybar, roxyos-hyprlock
- roxyos-kitty, roxyos-rofi, roxyos-plymouth
- roxyos-assets, roxyos-fish, roxyos-mako
- roxyos-sddm, roxyos-ly

## License

RNC-1 License ([LICENSE](LICENSE))

---

Made with my love for Roxy
