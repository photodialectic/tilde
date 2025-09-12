# ~
My dotfiles repository: modern Vim, tmux, and screen configurations with an enhanced installer.

## Features

- **Modern Vim Setup**: Migrated to vim-plug and CoC.nvim for LSP-based code completion.
- **Enhanced Vim Plugins**: vim-go, vim-polyglot, lightline, bufferline, vim-fugitive, vim-gitgutter, surround, commentary, GitHub Copilot, and AI-assisted editing.
- **Modern tmux Configuration**: Ctrl‑a prefix, 256‑color terminals, mouse support, vim-like pane navigation, pane resizing, and optional TPM plugin manager.
- **Enhanced screen Configuration**: Improved scrollback, UTF‑8 support, splits, vim‑like navigation, and a clean status line.
- **Go-based Installer**: Fast installation tool written in Go with better error handling, cross‑platform support, and automatic backups.
- **Smart Conflict Resolution**: Existing files are automatically backed up before being replaced by symlinks.

## Installation

### Quick Install (Recommended)

```bash
cd ~
git clone https://github.com/janmoesen/tilde.git
cd tilde
./install.sh
```

### Go-based Installer (Enhanced)

If you have Go installed, you can use the Go installer:

```bash
go build install.go
./install --prefix="$HOME" [--dry-run] [--verbose]
```

### Tilde Container CLI (like Claudex)

Build a standalone CLI that embeds the Dockerfile and dotfiles, then run a container mounting any directories you pass (similar to `claudex`).

Build the CLI binary (requires Go 1.19+):

```bash
cd tilde
go build -tags tildecli -o tildevim
```

Usage examples:

```bash
# Build/update the Docker image from embedded context
./tildevim build

# Build image and pre-install plugins at build time (requires network)
./tildevim build --with-plugins

# Mount current directory (as /work/<basename>) and open a bash shell
./tildevim

# Mount multiple projects, open bash
./tildevim ~/proj1 ~/proj2

# Mount a single project and open vim directly
./tildevim ~/proj --vim

# Open vim with a specific file
./tildevim ~/proj --vim README.md

# Pass environment variables
./tildevim ~/proj --env API_KEY=abc123 --env DEBUG

# Load variables from an env file
./tildevim ~/proj --env-file ~/proj/.env

# Auto-detect .env files in mounted dirs
./tildevim ~/proj1 ~/proj2 --auto-env
```

Notes:
- The container is named `tilde` and the image tag is `tilde`.
- Mounts appear under `/work/<basename>` inside the container.
- On first run, a Git repo is initialized in `/work` (branch `main`).
- `--env VAR[=VAL]` repeats and passes host env when VAL omitted.
- `--env-file PATH` repeats and passes each file to Docker `--env-file`.
- `--auto-env` discovers `.env` in each mounted directory and loads them.

### One-Step Build (CLI + Image)

If you want the Go build to also build the Docker image, use the provided Makefile:

```bash
cd tilde
make            # builds the CLI as ./tilde and then builds the Docker image via ./tilde build
# or
make cli        # builds just the CLI
make image      # builds the CLI then builds the Docker image
./tilde build --with-plugins  # optional: pre-install plugins during image build
```

The CLI also builds the image on first run if it's missing. If plugins are not present at runtime, it will attempt a best-effort `PlugInstall` on first start.

## Manual Plugin Setup

### Vim Plugins

Vim plugins are installed automatically by `install.sh`. To install or update plugins manually:

```bash
vim +PlugInstall +qall
```
