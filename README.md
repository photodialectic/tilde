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

## Manual Plugin Setup

### Vim Plugins

Vim plugins are installed automatically by `install.sh`. To install or update plugins manually:

```bash
vim +PlugInstall +qall
```
