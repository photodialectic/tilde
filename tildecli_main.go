//go:build tildecli
// +build tildecli

package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// We embed the minimal build context needed by the Dockerfile.
//
//go:embed Dockerfile .vimrc .tmux.conf .screenrc install-vim-plug.sh
var dockerContextFS embed.FS

func usage(prog string) {
	fmt.Printf(`Usage: %s [options] [DIR1 DIR2 ...]

Mounts each DIR at /work/<basename(DIRi)> in the tilde container.
If no DIR is provided, mounts the current directory at /work/<basename($PWD)>.

Options:
  build           Build or update the tilde Docker image from embedded context
  --vim[=bool]    Open vim after starting (default true; pass --vim=false for a shell)
  --env VAR[=VAL] Pass an environment variable (repeatable)
  --env-file PATH Load variables from an env file (repeatable)
  --auto-env      Auto-load .env files from mounted dirs
  --no-env-all    Do not forward all host environment variables
  --env-exclude K Exclude env var K when forwarding all (repeatable)
  --no-config-mount  Do not mount host config directory (~/.config)
  --config-path PATH  Override host config directory to mount

Examples:
  %s build
  %s .
  %s ~/proj1 ~/proj2 --vim
  %s ~/proj --vim README.md
`, prog, prog, prog, prog, prog)
}

func prepareBuildContext() (string, error) {
	tmpDir, err := os.MkdirTemp("", "tilde-build-")
	if err != nil {
		return "", err
	}
	// Write each embedded file into the temp dir
	entries := []string{"Dockerfile", ".vimrc", ".tmux.conf", ".screenrc", "install-vim-plug.sh"}
	for _, name := range entries {
		data, err := dockerContextFS.ReadFile(name)
		if err != nil {
			os.RemoveAll(tmpDir)
			return "", fmt.Errorf("cannot read embedded %s: %w", name, err)
		}
		out := filepath.Join(tmpDir, name)
		if err := os.WriteFile(out, data, 0644); err != nil {
			os.RemoveAll(tmpDir)
			return "", fmt.Errorf("cannot write %s: %w", out, err)
		}
	}
	return tmpDir, nil
}

func buildImage(withPlugins bool) error {
	fmt.Println("Building/updating the 'tilde' container image...")
	ctxDir, err := prepareBuildContext()
	if err != nil {
		return err
	}
	defer os.RemoveAll(ctxDir)
	args := []string{"build", "-t", "tilde"}
	if withPlugins {
		args = append(args, "--build-arg", "INSTALL_PLUGINS=1")
	}
	args = append(args, ctxDir)
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ensureImage() error {
	out, err := exec.Command("docker", "images", "-q", "tilde").Output()
	if err != nil {
		return fmt.Errorf("docker images check failed: %w", err)
	}
	if len(bytes.TrimSpace(out)) == 0 {
		// Default: install plugins at build time for faster container startup
		return buildImage(true)
	}
	return nil
}

func run(args []string) error {
	// Flags
	prog := filepath.Base(os.Args[0])
	fs := flag.NewFlagSet(prog, flag.ContinueOnError)
	// Default behavior now opens vim, mimicking host 'vim' usage
	openVim := fs.Bool("vim", true, "Open vim after starting (optionally with a file)")
	// Env passing
	var envs multiFlag
	var envFiles multiFlag
	autoEnv := fs.Bool("auto-env", false, "Auto-load .env files from mounted dirs")
	// Forward all host env vars by default; allow opting out and excluding names
	noEnvAll := fs.Bool("no-env-all", false, "Do not forward all host environment variables")
	var envExclude multiFlag
	fs.Var(&envs, "env", "Environment variable to pass (repeatable)")
	fs.Var(&envFiles, "env-file", "Env file to load (repeatable)")
	fs.Var(&envExclude, "env-exclude", "Environment variable name to exclude when forwarding all (repeatable)")

	// Config mounting: mount host config dir to container user's config path by default
	noConfigMount := fs.Bool("no-config-mount", false, "Do not mount host config directory (~/.config) into the container")
	configPath := fs.String("config-path", "", "Host config directory to mount (default: your XDG config dir)")
	// We need to allow an optional vim target argument; we parse leftover args later
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			usage(prog)
			return nil
		}
		return err
	}

	leftover := fs.Args()

	// Split mounts vs optional vim target
	var dirs []string
	var vimTarget string
	for i, a := range leftover {
		if *openVim && i == len(leftover)-1 && !isDir(a) {
			// last arg looks like a file; treat it as vim target
			vimTarget = a
			break
		}
		dirs = append(dirs, a)
	}

	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	if err := ensureImage(); err != nil {
		return err
	}

	// Remove any existing 'tilde' container
	_ = exec.Command("docker", "rm", "-f", "tilde").Run()

	// Build volume mounts
	var mounts []string
	for _, d := range dirs {
		abs, err := filepath.Abs(d)
		if err != nil {
			return fmt.Errorf("invalid path: %s", d)
		}
		fi, err := os.Stat(abs)
		if err != nil || !fi.IsDir() {
			return fmt.Errorf("'%s' is not a directory", abs)
		}
		name := filepath.Base(abs)
		mounts = append(mounts, "-v", fmt.Sprintf("%s:/work/%s", abs, name))
	}

	// Mount host config directory by default to mimic host Vim behavior
	if !*noConfigMount {
		hostCfg := *configPath
		if hostCfg == "" {
			if xdg, err := os.UserConfigDir(); err == nil && xdg != "" {
				hostCfg = xdg
			} else if home, err := os.UserHomeDir(); err == nil && home != "" {
				hostCfg = filepath.Join(home, ".config")
			}
		}
		if hostCfg != "" {
			if fi, err := os.Stat(hostCfg); err == nil && fi.IsDir() {
				mounts = append(mounts, "-v", fmt.Sprintf("%s:/home/node/.config", hostCfg))
			}
		}
	}

	// Auto-detect .env files if requested
	if *autoEnv {
		for _, d := range dirs {
			cand := filepath.Join(mustAbs(d), ".env")
			if fi, err := os.Stat(cand); err == nil && !fi.IsDir() {
				envFiles = append(envFiles, cand)
			}
		}
	}

	runArgs := []string{"run", "--name", "tilde", "-d"}
	runArgs = append(runArgs, mounts...)
	// Apply environment variables
	// 1) Forward all host env vars by default, with a safe denylist
	if !*noEnvAll {
		deny := map[string]struct{}{
			// Avoid breaking container runtime assumptions
			"PATH":   {},
			"PWD":    {},
			"OLDPWD": {},
			"SHLVL":  {},
			"_":      {},
			"HOME":   {},
		}
		// User-provided excludes
		for _, k := range envExclude {
			deny[k] = struct{}{}
		}
		for _, kv := range os.Environ() {
			// SplitN to allow empty values
			parts := strings.SplitN(kv, "=", 2)
			key := parts[0]
			if _, skip := deny[key]; skip {
				continue
			}
			// Keep as explicit KEY=VAL to avoid relying on caller env at docker invocation time
			runArgs = append(runArgs, "-e", kv)
		}
	}

	// 2) Apply explicitly requested envs (overrides)
	for _, e := range envs {
		// Allow either VAR or VAR=VAL
		if strings.Contains(e, "=") {
			runArgs = append(runArgs, "-e", e)
		} else {
			// Let docker pick from host env by name
			runArgs = append(runArgs, "-e", e)
		}
	}
	// 3) Apply env-files
	for _, f := range envFiles {
		af := f
		if !filepath.IsAbs(af) {
			af = mustAbs(af)
		}
		if fi, err := os.Stat(af); err == nil && !fi.IsDir() {
			runArgs = append(runArgs, "--env-file", af)
		} else {
			return fmt.Errorf("env-file not found: %s", f)
		}
	}
	runArgs = append(runArgs, "tilde", "sleep", "infinity")

	cmdRun := exec.Command("docker", runArgs...)
	cmdRun.Stdout = os.Stdout
	cmdRun.Stderr = os.Stderr
	if err := cmdRun.Run(); err != nil {
		return fmt.Errorf("docker run failed: %w", err)
	}

	// If plugins are not installed yet, attempt a best-effort install once
	// This avoids build-time network issues breaking the image
	pluginInit := `set -e
export HOME=/home/node
if [ ! -d "$HOME/.vim/plugged" ] || [ -z "$(ls -A "$HOME/.vim/plugged" 2>/dev/null)" ]; then
  echo 'Installing Vim plugins (first run)...'
  vim -E -s -u "$HOME/.vimrc" +'PlugInstall --sync' +qall || true
fi`
	cmdPlug := exec.Command("docker", "exec", "tilde", "bash", "-lc", pluginInit)
	cmdPlug.Stdout = os.Stdout
	cmdPlug.Stderr = os.Stderr
	_ = cmdPlug.Run()

	// Decide what to attach to
	if *openVim {
		// If a single directory was provided, start in that subdir; otherwise /work
		startDir := "/work"
		if len(dirs) == 1 {
			base := filepath.Base(mustAbs(dirs[0]))
			startDir = filepath.Join("/work", base)
		}
		var cmd *exec.Cmd
		if vimTarget != "" {
			// Quote target lightly; rely on docker exec args
			cmd = exec.Command("docker", "exec", "-it", "tilde", "bash", "-lc", fmt.Sprintf("cd %s && vim %s", shellQuote(startDir), shellQuote(vimTarget)))
		} else {
			cmd = exec.Command("docker", "exec", "-it", "tilde", "bash", "-lc", fmt.Sprintf("cd %s && vim", shellQuote(startDir)))
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Default: open an interactive bash shell
	cmd := exec.Command("docker", "exec", "-it", "tilde", "bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isDir(p string) bool {
	abs, err := filepath.Abs(p)
	if err != nil {
		return false
	}
	fi, err := os.Stat(abs)
	return err == nil && fi.IsDir()
}

func mustAbs(p string) string {
	a, _ := filepath.Abs(p)
	return a
}

func shellQuote(s string) string {
	if s == "" {
		return ""
	}
	// crude escaping for spaces and quotes
	if strings.IndexFunc(s, func(r rune) bool { return r == ' ' || r == '\'' || r == '"' }) == -1 {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

// multiFlag collects repeated flag values
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "build" {
			// Support build options, e.g., --with-plugins
			bfs := flag.NewFlagSet("build", flag.ContinueOnError)
			withPlugins := bfs.Bool("with-plugins", true, "Pre-install Vim plugins at image build time")
			noPlugins := bfs.Bool("no-plugins", false, "Skip plugin installation at build time")
			_ = bfs.Parse(os.Args[2:])
			installPlugins := *withPlugins && !*noPlugins
			if err := buildImage(installPlugins); err != nil {
				fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
				os.Exit(1)
			}
			return
		}
		if os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help" {
			usage(filepath.Base(os.Args[0]))
			return
		}
	}
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
