BIN?=tilde
INSTALL_DIR?=/usr/local/bin

.PHONY: all build image clean install uninstall

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
install: build
	@echo "Installing $(BIN) to $(INSTALL_DIR)"
	@if [ -w "$(INSTALL_DIR)" ]; then \
		cp $(BIN) $(INSTALL_DIR)/; \
	else \
		sudo cp $(BIN) $(INSTALL_DIR)/; \
	fi
	@echo "✅ Installed $(BIN) to $(INSTALL_DIR)"

# Uninstall binary from system directory
uninstall:
	@echo "Removing $(BIN) from $(INSTALL_DIR)"
	@if [ -w "$(INSTALL_DIR)" ]; then \
		rm -f $(INSTALL_DIR)/$(BIN); \
	else \
		sudo rm -f $(INSTALL_DIR)/$(BIN); \
	fi
	@echo "✅ Removed $(BIN) from $(INSTALL_DIR)"

# Clean up build artifacts
clean:
	rm -f $(BIN)
	docker rmi tilde 2>/dev/null || true
