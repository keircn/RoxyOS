set -g roxyos_color_primary "#1a5276"
set -g roxyos_color_accent "#5dade2"
set -g roxyos_color_background "#0a0a0f"
set -g roxyos_color_text "#e8f4fc"

set -g fish_color_normal $roxyos_color_text
set -g fish_color_command $roxyos_color_accent
set -g fish_color_param "#85c1e9"
set -g fish_color_keyword $roxyos_color_primary
set -g fish_color_quote "#d4c9a8"
set -g fish_color_redirection "#a8c8d8"
set -g fish_color_end "#4a6b5a"
set -g fish_color_error "#8b4a2b"
set -g fish_color_gray "#5dade2"
set -g fish_color_selection --background=$roxyos_color_primary
set -g fish_color_search_match --background=$roxyos_color_accent
set -g fish_color_history_current $roxyos_color_text
set -g fish_color_operator "#85c1e9"
set -g fish_color_escape "#c9a227"
set -g fish_color_cwd $roxyos_color_accent
set -g fish_color_cwd_root "#8b4a2b"

set -g fish_pager_color_prefix $roxyos_color_text
set -g fish_pager_color_completion $roxyos_color_text
set -g fish_pager_color_description "#a8c8d8"
set -g fish_pager_color_progress $roxyos_color_accent
