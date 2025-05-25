package main

import (
	"bufio"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	DryRun    bool
	TargetDir string
	SourceDir string
	Verbose   bool
}

var ignoredFiles = []string{
	".gitignore",
	".gitmodules",
	".git",
	"COPYING",
	"README.md",
	"install.sh",
	"install.go",
	"go.mod",
	"go.sum",
}

func main() {
	config := &Config{}

	flag.BoolVar(&config.DryRun, "dry-run", false, "Show what would be done without executing")
	flag.StringVar(&config.TargetDir, "prefix", os.Getenv("HOME"), "Target directory for installation")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if config.DryRun {
		fmt.Println("🚀 Dry run mode. Not actually executing anything.")
	}

	// Get absolute paths
	var err error
	config.SourceDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Failed to get source directory: %v", err)
	}

	config.TargetDir, err = filepath.Abs(config.TargetDir)
	if err != nil {
		log.Fatalf("Failed to get target directory: %v", err)
	}

	if err := validateTargetDir(config.TargetDir); err != nil {
		log.Fatalf("Target directory validation failed: %v", err)
	}

	fmt.Printf("📂 Installing from %s to %s\n", config.SourceDir, config.TargetDir)

	// Install vim-plug if not present
	if err := installVimPlug(config); err != nil {
		log.Printf("Warning: Failed to install vim-plug: %v", err)
	}

	// Get list of files to symlink
	sourceFiles, err := getSourceFiles(config.SourceDir)
	if err != nil {
		log.Fatalf("Failed to get source files: %v", err)
	}

	// Check for conflicts
	conflicts, err := findConflicts(config, sourceFiles)
	if err != nil {
		log.Fatalf("Failed to check for conflicts: %v", err)
	}

	if len(conflicts) > 0 {
		if !handleConflicts(config, conflicts) {
			fmt.Println("❌ Installation aborted by user.")
			os.Exit(1)
		}
	}

	// Create symlinks
	if err := createSymlinks(config, sourceFiles); err != nil {
		log.Fatalf("Failed to create symlinks: %v", err)
	}

	// Install vim plugins
	if err := installVimPlugins(config); err != nil {
		log.Printf("Warning: Failed to install vim plugins: %v", err)
	}

	fmt.Println("✅ Installation complete!")
}

func validateTargetDir(targetDir string) error {
	if targetDir == "" {
		return fmt.Errorf("target directory is empty")
	}
	if targetDir == "/" {
		return fmt.Errorf("refusing to install into root directory")
	}

	info, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		return os.MkdirAll(targetDir, 0755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", targetDir)
	}

	// Check if writable
	testFile := filepath.Join(targetDir, ".tilde_write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("%s is not writable", targetDir)
	}
	f.Close()
	os.Remove(testFile)

	return nil
}

func installVimPlug(config *Config) error {
	vimPlugPath := filepath.Join(config.TargetDir, ".vim/autoload/plug.vim")

	if _, err := os.Stat(vimPlugPath); err == nil {
		fmt.Println("📦 vim-plug already installed")
		return nil
	}

	fmt.Println("📦 Installing vim-plug...")

	if config.DryRun {
		fmt.Printf("Would download vim-plug to %s\n", vimPlugPath)
		return nil
	}

	// Create directory
	if err := os.MkdirAll(filepath.Dir(vimPlugPath), 0755); err != nil {
		return err
	}

	// Try curl first, then wget
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-Command",
			fmt.Sprintf("Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim' -OutFile '%s'", vimPlugPath))
	default:
		// Try curl first
		cmd = exec.Command("curl", "-fLo", vimPlugPath, "--create-dirs",
			"https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim")
		if err := cmd.Run(); err != nil {
			// Fall back to wget
			cmd = exec.Command("wget", "-O", vimPlugPath,
				"https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to download vim-plug with both curl and wget: %v", err)
			}
		}
		fmt.Println("✅ vim-plug installed successfully")
		return nil
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download vim-plug: %v. Try running ./install-vim-plug.sh manually when network is available", err)
	}

	fmt.Println("✅ vim-plug installed successfully")
	return nil
}

func getSourceFiles(sourceDir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Skip ignored files
		for _, ignored := range ignoredFiles {
			if strings.Contains(relPath, ignored) {
				return nil
			}
		}

		files = append(files, relPath)
		return nil
	})

	return files, err
}

func findConflicts(config *Config, sourceFiles []string) ([]string, error) {
	var conflicts []string

	for _, file := range sourceFiles {
		targetPath := filepath.Join(config.TargetDir, file)
		sourcePath := filepath.Join(config.SourceDir, file)

		// Check if target exists
		if _, err := os.Lstat(targetPath); os.IsNotExist(err) {
			continue
		}

		// Check if they're the same file (already symlinked)
		if os.SameFile != nil {
			sourceInfo, err1 := os.Stat(sourcePath)
			targetInfo, err2 := os.Stat(targetPath)
			if err1 == nil && err2 == nil && os.SameFile(sourceInfo, targetInfo) {
				continue
			}
		}

		// Check if contents are the same
		if filesEqual(sourcePath, targetPath) {
			continue
		}

		conflicts = append(conflicts, file)
	}

	return conflicts, nil
}

func filesEqual(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		return false
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false
	}
	defer f2.Close()

	h1 := md5.New()
	h2 := md5.New()

	_, err1 := io.Copy(h1, f1)
	_, err2 := io.Copy(h2, f2)

	if err1 != nil || err2 != nil {
		return false
	}

	return string(h1.Sum(nil)) == string(h2.Sum(nil))
}

func handleConflicts(config *Config, conflicts []string) bool {
	backupSuffix := fmt.Sprintf(".tilde-%s", time.Now().Format("20060102-150405"))

	fmt.Printf("⚠️  WARNING: %d file(s) conflict with Tilde.\n", len(conflicts))
	fmt.Printf("Your files will be backed up with suffix: %s\n\n", backupSuffix)

	for _, file := range conflicts {
		targetPath := filepath.Join(config.TargetDir, file)
		fmt.Printf("  📄 %s\n", targetPath)
	}

	fmt.Print("\n🤔 Continue with installation? (yes/no): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "yes" {
		return false
	}

	// Backup conflicting files
	for _, file := range conflicts {
		targetPath := filepath.Join(config.TargetDir, file)
		backupPath := targetPath + backupSuffix

		if config.DryRun {
			fmt.Printf("Would backup: %s -> %s\n", targetPath, backupPath)
		} else {
			if err := os.Rename(targetPath, backupPath); err != nil {
				log.Printf("Failed to backup %s: %v", targetPath, err)
			} else {
				fmt.Printf("📦 Backed up: %s\n", file)
			}
		}
	}

	return true
}

func createSymlinks(config *Config, sourceFiles []string) error {
	linkCount := 0

	for _, file := range sourceFiles {
		sourcePath := filepath.Join(config.SourceDir, file)
		targetPath := filepath.Join(config.TargetDir, file)

		// Create target directory if needed
		targetDir := filepath.Dir(targetPath)
		if config.DryRun {
			fmt.Printf("Would create directory: %s\n", targetDir)
		} else {
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
			}
		}

		// Check if symlink already exists
		if linkTarget, err := os.Readlink(targetPath); err == nil {
			if linkTarget == sourcePath {
				continue // Already correctly symlinked
			}
		}

		// Create relative symlink
		relSource, err := filepath.Rel(targetDir, sourcePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %v", err)
		}

		if config.DryRun {
			fmt.Printf("Would symlink: %s -> %s\n", targetPath, relSource)
		} else {
			if err := os.Symlink(relSource, targetPath); err != nil {
				return fmt.Errorf("failed to create symlink %s: %v", targetPath, err)
			}
			if config.Verbose {
				fmt.Printf("🔗 Linked: %s\n", file)
			}
			linkCount++
		}
	}

	if linkCount == 0 && !config.DryRun {
		fmt.Println("📋 All files were already symlinked.")
	} else if !config.DryRun {
		fmt.Printf("🔗 Created %d symlinks\n", linkCount)
	}

	return nil
}

func installVimPlugins(config *Config) error {
	fmt.Println("📦 Installing vim plugins...")

	if config.DryRun {
		fmt.Println("Would run: vim +PlugInstall +qall")
		return nil
	}

	cmd := exec.Command("vim", "+PlugInstall", "+qall")
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("vim plugin installation failed: %v", err)
	}

	fmt.Println("✅ Vim plugins installed")
	return nil
}
