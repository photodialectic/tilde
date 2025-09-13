BIN?=tilde
INSTALL_DIR?=/usr/local/bin

.PHONY: all build image clean install install-cli install-dotfiles uninstall uninstall-dotfiles

all: build image

# Build the CLI binary with the tildecli tag
build:
	go build -tags tildecli -o $(BIN)

# Build/update the Docker image via the CLI's embedded context
image: build
	./$(BIN) build

# Force a rebuild of the Docker image directly (bypass CLI)
rebuild-image:
	docker build -t tilde .

# Install binary to system directory
install-cli: build
	@echo "Installing $(BIN) to $(INSTALL_DIR)"
	@if [ -w "$(INSTALL_DIR)" ]; then \
		cp $(BIN) $(INSTALL_DIR)/; \
	else \
		sudo cp $(BIN) $(INSTALL_DIR)/; \
	fi
	@echo "✅ Installed $(BIN) to $(INSTALL_DIR)"

# Install dotfiles via symlinks using Go installer
install-dotfiles:
	@echo "Installing dotfiles via symlinks"
	@if [ -f dotfiles.go ] && command -v go >/dev/null 2>&1; then \
		go run dotfiles.go; \
	else \
		echo "❌ Go not found - required for dotfile installation"; \
		exit 1; \
	fi
	@echo "✅ Dotfiles installed"

# Uninstall dotfiles and restore backups
uninstall-dotfiles:
	@echo "Uninstalling dotfiles and restoring backups"
	@if [ -f dotfiles.go ] && command -v go >/dev/null 2>&1; then \
		go run dotfiles.go --uninstall; \
	else \
		echo "❌ Go not found - required for dotfile uninstallation"; \
		exit 1; \
	fi
	@echo "✅ Dotfiles uninstalled"

# Install both CLI and dotfiles
install: install-cli install-dotfiles

# Uninstall binary from system directory
uninstall-cli:
	@echo "Removing $(BIN) from $(INSTALL_DIR)"
	@if [ -w "$(INSTALL_DIR)" ]; then \
		rm -f $(INSTALL_DIR)/$(BIN); \
	else \
		sudo rm -f $(INSTALL_DIR)/$(BIN); \
	fi
	@echo "✅ Removed $(BIN) from $(INSTALL_DIR)"

# Uninstall both CLI and dotfiles
uninstall: uninstall-cli uninstall-dotfiles

# Clean up build artifacts
clean:
	rm -f $(BIN)
	docker rmi tilde 2>/dev/null || true
