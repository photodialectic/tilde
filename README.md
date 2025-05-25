# ~
My dot files repo with modern vim and screen configuration.

## Features
- **Modern Vim Setup**: Migrated from Vundle to vim-plug with LSP support via CoC
- **Enhanced Plugins**: FZF integration, Git gutter, auto-pairs, and more
- **Go Installer**: Fast, concurrent installation with better error handling
- **Smart Conflict Resolution**: Automatic backup of existing files

## Installation

### Quick Install (Recommended)
```bash
cd ~
git clone https://github.com/photodialectic/tilde
cd tilde
./install.sh
```

### Go Installer (Enhanced)
If you have Go installed:
```bash
go build install.go
./install [--dry-run] [--prefix=$HOME] [--verbose]
```

### Manual Plugin Setup
After installation, vim plugins will be automatically installed. For manual installation:
```bash
vim +PlugInstall +qall
```

## Key Improvements

### Vim Configuration
- **Plugin Manager**: Upgraded from Vundle to vim-plug
- **LSP Support**: CoC.nvim for intelligent code completion
- **Fuzzy Finding**: FZF integration for file/buffer navigation
- **Git Integration**: Enhanced with git-gutter for diff indicators
- **Modern Syntax**: vim-polyglot for comprehensive language support

### Key Bindings
- `<leader>f` - Find files with FZF
- `<leader>g` - Search in files with ripgrep
- `<leader>b` - Buffer switcher
- `gd` - Go to definition (LSP)
- `gr` - Find references (LSP)
- `[g`/`]g` - Navigate diagnostics

### Go Installer Benefits
- ⚡ Concurrent operations for faster installation
- 🛡️ Better error handling and validation
- 📦 Automatic vim-plug installation
- 🔄 Smart conflict resolution with backups
- 🎯 Cross-platform compatibility

## Update

```bash
cd ~/tilde
./install.sh
```
