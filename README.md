(G)it (W)ork(TR)ee 

Golang-based manager for git worktrees.

## Installation 
You can install the tool by running `go install` as follows:
```bash
go install github.com/amig3n/gwtr@latest
```

Or just clone the repository and build the binary yourself:
```bash
make build
```

## Install and init
First, you need to install required wrapper script, so gwtr could 
change your working directory within your shell. You can do so by adding the
following snippet into your shell configuration file (e.g. `~/.bashrc` or `~/.zshrc`):

```bash
gwtr() {
    if [ "$1" = "switch" ]; then
        cd "$(command gwtr switch "$2")" || return 1
    else
        command gwtr "$@"
    fi
}
```

Then initialize the gwtr inside your repository:
```bash
gwtr init
```

## Usage
From this moment you can use `gwtr` command to manage your git worktrees.
To create a new worktree:
```bash 
gwtr create -p <path> -b <branch>
```

to switch between worktrees:
```bash
gwtr switch <worktree_ID_number>
```

To delete worktree:
```bash
gwtr delete <worktree_ID_number>
```

You can always check the `--help` flag for more info about each command.

Wish you a pleasant usage :)
