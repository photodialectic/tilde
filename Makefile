BIN?=tilde

.PHONY: all cli image rebuild-image clean

all: cli image

# Build the CLI binary with the tildecli tag
cli:
	go build -tags tildecli -o $(BIN)

# Build/update the Docker image via the CLI's embedded context
image: cli
	./$(BIN) build

# Force a rebuild of the Docker image directly (bypass CLI)
rebuild-image:
	docker build -t tilde .

clean:
	rm -f $(BIN)
