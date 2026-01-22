.PHONY: build install uninstall clean test

BINARY_NAME=vibe-monitor
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=.
INSTALL_PATH?=/usr/local/bin

# Detect user's local bin if it exists
LOCAL_BIN=$(HOME)/.local/bin
ifeq ($(shell test -d $(LOCAL_BIN) && echo yes),yes)
    DEFAULT_INSTALL_PATH=$(LOCAL_BIN)
else
    DEFAULT_INSTALL_PATH=/usr/local/bin
endif

build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@if [ "$(INSTALL_PATH)" = "/usr/local/bin" ] && [ ! -w "$(INSTALL_PATH)" ]; then \
		echo "Installing to system path requires sudo..."; \
		sudo install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME); \
	else \
		install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME); \
	fi
	@echo "✓ Installed to $(INSTALL_PATH)/$(BINARY_NAME)"
	@echo ""
	@echo "To verify installation, run: $(BINARY_NAME) --version"

install-user: INSTALL_PATH=$(LOCAL_BIN)
install-user: install
	@echo "Note: Make sure $(LOCAL_BIN) is in your PATH"

uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_PATH)..."
	@if [ "$(INSTALL_PATH)" = "/usr/local/bin" ] && [ ! -w "$(INSTALL_PATH)" ]; then \
		sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	else \
		rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	fi
	@echo "✓ Uninstalled $(BINARY_NAME)"

clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@echo "✓ Clean complete"

test:
	@echo "Running tests..."
	go test -v ./...

test-fastfetch: build
	@echo "Testing fastfetch integration..."
	@./$(BINARY_NAME) --test-fastfetch

help:
	@echo "vibe-monitor Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build          Build the binary"
	@echo "  make install        Install to $(DEFAULT_INSTALL_PATH) (system-wide)"
	@echo "  make install-user   Install to ~/.local/bin (user only)"
	@echo "  make uninstall      Remove installed binary"
	@echo "  make clean          Remove build artifacts"
	@echo "  make test           Run tests"
	@echo "  make test-fastfetch Test fastfetch integration"
	@echo ""
	@echo "Variables:"
	@echo "  INSTALL_PATH        Override installation path (default: $(DEFAULT_INSTALL_PATH))"
	@echo ""
	@echo "Examples:"
	@echo "  make install-user                      # Install to ~/.local/bin"
	@echo "  make INSTALL_PATH=~/bin install        # Install to custom location"
