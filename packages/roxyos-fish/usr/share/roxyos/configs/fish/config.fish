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

set -g fish_greeting

set fish_pager_color_prefix cyan
set fish_color_autosuggestion brblack
set EDITOR nvim

abbr n nvim

abbr .. 'cd ..'
abbr ... 'cd ../..'
abbr .3 'cd ../../..'
abbr .4 'cd ../../../..'
abbr .5 'cd ../../../../..'

abbr g git
abbr ga 'git add'
abbr gb 'git branch'
abbr gc 'git commit'
abbr gca 'git commit --amend'
abbr gco 'git checkout'
abbr gd 'git diff'
abbr gl 'git pull'
abbr gp 'git push'
abbr gst 'git status'
abbr grv 'git remote -v'
abbr glg 'git log --oneline --graph --decorate --all'
abbr gci 'git commit -a -m "Initial commit"'

abbr mkdir 'mkdir -p'
abbr df 'df -h'
abbr du 'du -h --max-depth=1'
abbr free 'free -h'
abbr pls sudo

fish_vi_key_bindings

if command -q starship
    starship init fish | source
end

function roxyos_welcome
    set -l cols (tput cols)
    set -l box "╭────────────────────────────────────╮"
    set -l mid "│       Welcome to RoxyOS            │"
    set -l bot "╰────────────────────────────────────╯"
    set -l box_width 38
    set -l pad (math "($cols - $box_width) / 2")
    set -l padding (string repeat -n $pad " ")
    echo ""
    set_color 5dade2
    echo "$padding$box"
    echo "$padding$mid"
    echo "$padding$bot"
    set_color normal
    echo ""
end

if status is-login; and test -z "$ROXYOS_WELCOMED"
    set -gx ROXYOS_WELCOMED 1
    roxyos_welcome
end

if test -f $HOME/.config/fish/conf.d/custom.fish
    source $HOME/.config/fish/conf.d/custom.fish
end
