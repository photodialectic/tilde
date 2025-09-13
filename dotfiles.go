//go:build !tildecli
// +build !tildecli

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Config struct {
	DryRun    bool
	TargetDir string
	SourceDir string
	Uninstall bool
	BackupDir string
}

var dotfiles = []string{
	".vimrc",
	".tmux.conf",
	".screenrc",
}

func main() {
	config := &Config{}

	flag.BoolVar(&config.DryRun, "dry-run", false, "Show what would be done without executing")
	flag.StringVar(&config.TargetDir, "prefix", os.Getenv("HOME"), "Target directory for installation")
	flag.BoolVar(&config.Uninstall, "uninstall", false, "Uninstall dotfiles and restore backups")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// Set source directory to current working directory (where the repo is)
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}
	config.SourceDir = workingDir

	// Set backup directory
	config.BackupDir = filepath.Join(config.TargetDir, ".tilde-backup")

	if config.Uninstall {
		if err := uninstallDotfiles(config); err != nil {
			fmt.Printf("❌ Uninstall failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Tilde dotfiles uninstalled successfully")
	} else {
		if err := installDotfiles(config); err != nil {
			fmt.Printf("❌ Install failed: %v\n", err)
			os.Exit(1)
		}

		if err := installVimPlugins(config); err != nil {
			fmt.Printf("⚠️  Vim plugins installation failed: %v\n", err)
		}

		fmt.Println("✅ Tilde installation completed successfully")
	}
}

func installDotfiles(config *Config) error {
	// Ensure backup directory exists
	if !config.DryRun {
		if err := os.MkdirAll(config.BackupDir, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}
	}

	for _, dotfile := range dotfiles {
		sourcePath := filepath.Join(config.SourceDir, dotfile)
		targetPath := filepath.Join(config.TargetDir, dotfile)
		backupPath := filepath.Join(config.BackupDir, dotfile)

		// Check if source file exists
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			fmt.Printf("⚠️  Skipping %s (not found in tilde repo)\n", dotfile)
			continue
		}

		// Handle existing target file
		if info, err := os.Lstat(targetPath); err == nil {
			fmt.Printf("🔍 Found existing %s at %s\n", dotfile, targetPath)
			if info.Mode()&os.ModeSymlink != 0 {
				// It's already a symlink - check if it points to our file
				if target, err := os.Readlink(targetPath); err == nil {
					resolvedTarget, err1 := filepath.EvalSymlinks(targetPath)
					resolvedSource, err2 := filepath.EvalSymlinks(sourcePath)
					if err1 == nil && err2 == nil && resolvedTarget == resolvedSource {
						fmt.Printf("✓ %s already correctly symlinked\n", dotfile)
						continue
					}
				}
				// Remove incorrect symlink
				fmt.Printf("🔄 Removing incorrect symlink: %s\n", dotfile)
				if !config.DryRun {
					if err := os.Remove(targetPath); err != nil {
						return fmt.Errorf("failed to remove incorrect symlink %s: %w", dotfile, err)
					}
				}
			} else {
				// It's a regular file - back it up
				fmt.Printf("💾 Backing up existing %s to %s\n", dotfile, backupPath)
				if !config.DryRun {
					if err := copyFile(targetPath, backupPath); err != nil {
						return fmt.Errorf("failed to backup %s: %w", dotfile, err)
					}
					if err := os.Remove(targetPath); err != nil {
						return fmt.Errorf("failed to remove %s after backup: %w", dotfile, err)
					}
				}
			}
		} else {
			fmt.Printf("📝 No existing %s found at %s\n", dotfile, targetPath)
		}

		// Create symlink using relative path for portability
		relSourcePath, err := filepath.Rel(filepath.Dir(targetPath), sourcePath)
		if err != nil {
			return fmt.Errorf("failed to compute relative path for %s: %w", dotfile, err)
		}
		fmt.Printf("🔗 Symlinking %s -> %s\n", dotfile, relSourcePath)
		if !config.DryRun {
			if err := os.Symlink(relSourcePath, targetPath); err != nil {
				return fmt.Errorf("failed to create symlink for %s: %w", dotfile, err)
			}
		}
	}

	return nil
}

func uninstallDotfiles(config *Config) error {
	for _, dotfile := range dotfiles {
		targetPath := filepath.Join(config.TargetDir, dotfile)
		backupPath := filepath.Join(config.BackupDir, dotfile)

		// Check if target is our symlink
		if info, err := os.Lstat(targetPath); err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				// Remove symlink
				fmt.Printf("🗑️  Removing symlink: %s\n", dotfile)
				if !config.DryRun {
					if err := os.Remove(targetPath); err != nil {
						fmt.Printf("⚠️  Failed to remove symlink %s: %v\n", dotfile, err)
						continue
					}
				}
			} else {
				fmt.Printf("⚠️  %s is not a symlink, skipping\n", dotfile)
				continue
			}
		} else {
			fmt.Printf("⚠️  %s not found, skipping\n", dotfile)
			continue
		}

		// Restore backup if it exists
		if _, err := os.Stat(backupPath); err == nil {
			fmt.Printf("🔄 Restoring backup: %s\n", dotfile)
			if !config.DryRun {
				if err := copyFile(backupPath, targetPath); err != nil {
					return fmt.Errorf("failed to restore backup for %s: %w", dotfile, err)
				}
				if err := os.Remove(backupPath); err != nil {
					fmt.Printf("⚠️  Failed to remove backup file %s: %v\n", backupPath, err)
				}
			}
		}
	}

	// Remove backup directory if empty
	if !config.DryRun {
		if entries, err := os.ReadDir(config.BackupDir); err == nil && len(entries) == 0 {
			if err := os.Remove(config.BackupDir); err != nil {
				fmt.Printf("⚠️  Failed to remove backup directory %s: %v\n", config.BackupDir, err)
			}
		}
	}

	return nil
}

func installVimPlugins(config *Config) error {
	if config.DryRun {
		fmt.Println("Would install vim plugins")
		return nil
	}

	// Check if vim is available
	if _, err := exec.LookPath("vim"); err != nil {
		return fmt.Errorf("vim not found in PATH")
	}

	fmt.Println("📦 Installing vim plugins...")

	// Install vim-plug if needed
	vimPlugPath := filepath.Join(config.TargetDir, ".vim", "autoload", "plug.vim")
	if _, err := os.Stat(vimPlugPath); os.IsNotExist(err) {
		fmt.Println("📥 Installing vim-plug...")
		if err := installVimPlug(config); err != nil {
			return err
		}
	}

	// Run vim plugin installation
	cmd := exec.Command("vim", "+PlugInstall", "+qall")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("vim plugin installation failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Println("✅ Vim plugins installed")
	return nil
}

func installVimPlug(config *Config) error {
	vimDir := filepath.Join(config.TargetDir, ".vim", "autoload")
	if err := os.MkdirAll(vimDir, 0755); err != nil {
		return fmt.Errorf("failed to create vim autoload directory: %w", err)
	}

	vimPlugPath := filepath.Join(vimDir, "plug.vim")
	cmd := exec.Command("curl", "-fLo", vimPlugPath,
		"https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim")

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to download vim-plug: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func copyFile(src, dst string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Copy permissions
	if info, err := srcFile.Stat(); err == nil {
		dstFile.Chmod(info.Mode())
	}

	return nil
}
