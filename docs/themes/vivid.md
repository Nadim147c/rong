# Vivid (`LS_COLORS`)

**vivid** is a generator for the **`LS_COLORS`** environment variable that controls
the colorized output of
[`ls`](https://www.gnu.org/software/coreutils/manual/html_node/ls-invocation.html#ls-invocation),
[`tree`](http://mama.indstate.edu/users/ice/tree/),
[`fd`](https://github.com/sharkdp/fd),
[`bfs`](https://github.com/tavianator/bfs),
[`dust`](https://github.com/bootandy/dust) and many other tools.

## Install

Install [vivid](https://github.com/sharkdp/vivid?tab=readme-ov-file#installation)
from your preferred package manager.

```bash
# With pacman
sudo pacman -S vivid
```

## Apply

Add following snippets to your preferred shell's init script.

- **Bash** or **ZSH**

```bash
export LS_COLORS=$(vivid generate "${XDG_STATE_HOME:-$HOME/.local/state}/rong/vivid.yml")
```

- **Fish**

```fish
set state_home (test -n "$XDG_STATE_HOME"; and echo $XDG_STATE_HOME; or echo "$HOME/.local/state")
set -x LS_COLORS (vivid generate "$state_home/rong/vivid.yml")
```

- **Nushell**

```nu
$env.LS_COLORS = (vivid generate $"($env.XDG_STATE_HOME | default $"($env.HOME)/.local/state")/rong/vivid.yml")
```

## Reload

Themes will reload on shell restart. Sourcing shell init script will also reload the
theme (`source ~/.bashrc`).

However, if live reload is a must than you can generate a `LS_COLORS` using
`dircolors`.

```bash
source <(dircolors)
```

This will use basic terminal colors and will be updated when you update the terminal.
