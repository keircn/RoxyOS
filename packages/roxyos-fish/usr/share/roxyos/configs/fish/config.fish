if test -f $HOME/.config/fish/env.fish
    source $HOME/.config/fish/env.fish
end

if test -f $HOME/.config/fish/conf.d/roxyos-theme.fish
    source $HOME/.config/fish/conf.d/roxyos-theme.fish
end

if test -f $HOME/.config/fish/conf.d/roxyos-env.fish
    source $HOME/.config/fish/conf.d/roxyos-env.fish
end

if test -f $HOME/.config/fish/conf.d/abbreviations.fish
    source $HOME/.config/fish/conf.d/abbreviations.fish
end

if command -q starship
    starship init fish | source
end

fish_vi_key_bindings

set -g fish_color_normal e8f4fc
set -g fish_color_command 5dade2
set -g fish_color_param 85c1e9
set -g fish_color_keyword 1a5276
set -g fish_color_quote d4c9a8
set -g fish_color_redirection a8c8d8
set -g fish_color_end 4a6b5a
set -g fish_color_error 8b4a2b
set -g fish_color_gray 5dade2
set -g fish_color_selection --background=1a5276
set -g fish_color_search_match --background=5dade2
set -g fish_color_history_current e8f4fc
set -g fish_color_operator 85c1e9
set -g fish_color_escape c9a227
set -g fish_color_cwd 5dade2
set -g fish_color_cwd_root 8b4a2b

set -g fish_pager_color_prefix e8f4fc
set -g fish_pager_color_completion e8f4fc
set -g fish_pager_color_description a8c8d8
set -g fish_pager_color_progress e8f4fc

if test -z "$ROXYOS_SKIP_WELCOME"
    echo ""
    echo -e "\033[38;5;74m╭────────────────────────────────────╮\033[0m"
    echo -e "\033[38;5;74m│\033[0m   \033[38;5;117mWelcome to RoxyOS\033[0m               \033[38;5;74m│\033[0m"
    echo -e "\033[38;5;74m╰────────────────────────────────────╯\033[0m"
    echo ""
end
