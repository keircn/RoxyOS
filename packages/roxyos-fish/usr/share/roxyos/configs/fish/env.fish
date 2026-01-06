set -gx XDG_CONFIG_HOME "$HOME/.config"
set -gx XDG_DATA_HOME "$HOME/.local/share"
set -gx XDG_STATE_HOME "$HOME/.local/state"
set -gx XDG_CACHE_HOME "$HOME/.cache"

set -gx DESKTOP_SESSION "roxyos-niri"
set -gx XDG_CURRENT_DESKTOP "niri"
set -gx XDG_SESSION_TYPE "wayland"

set -gx XCURSOR_SIZE "20"
set -gx XCURSOR_THEME "BreezeX-RoséPine-Linux"

set -gx GTK_THEME "Adwaita:dark"
set -gx GTK_ICON_THEME "Adwaita"
set -gx GTK_APPLICATION_PREFER_DARK_THEME "1"

set -gx QT_QPA_PLATFORMTHEME "qt6ct"
set -gx QT_STYLE_OVERRIDE "Adwaita-Dark"
set -gx QT_FONT_DPI "96"

set -gx MOZ_ENABLE_WAYLAND "1"
set -gx GDK_BACKEND "wayland,x11"
set -gx CLUTTER_BACKEND "wayland"
set -gx SDL_VIDEODRIVER "wayland"

set -gx _JAVA_AWT_WM_NONREPARENTING "1"

set -gx EDITOR "nvim"
set -gx VISUAL "nvim"
set -gx PAGER "less"
set -gx BROWSER "firefox"
set -gx TERMINAL "kitty"
set -gx SHELL "/usr/bin/fish"

set -gx PATH "$HOME/.local/bin" $PATH
set -gx PATH "$HOME/.cargo/bin" $PATH
set -gx PATH "$HOME/go/bin" $PATH

set -gx CARGO_HOME "$HOME/.cargo"
set -gx RUSTUP_HOME "$HOME/.rustup"
set -gx GOPATH "$HOME/go"

set -gx ROXYOS_THEME "migurdia"
set -gx ROXYOS_COLOR_PRIMARY "#1a5276"
set -gx ROXYOS_COLOR_ACCENT "#5dade2"
set -gx ROXYOS_COLOR_BACKGROUND "#0a1628"

set -gx TERM "xterm-256color"
