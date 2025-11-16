# --- Configuration ---
# Define the main application name.
# This corresponds to the directory in the 'cmd/' folder.
APP_NAME=webapi

# Define paths
BINARY_DIR=tmp
BINARY_PATH=$(BINARY_DIR)/$(APP_NAME)
MAIN_PKG=./cmd/$(APP_NAME)
TAILWIND_INPUT=./static/css/custom.css
TAILWIND_OUTPUT=./static/css/style.css

.PHONY: help
help: ## print make targets 
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# --- Tool Installation ---
.PHONY: go-install-air
go-install-air: ## Installs the air build reload system using 'go install'
	go install github.com/air-verse/air@latest

.PHONY: get-install-air
get-install-air: ## Installs the air build reload system using cUrl
	curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

.PHONY: go-install-templ
go-install-templ: ## Installs the templ Templating system for Go
	go install github.com/a-h/templ/cmd/templ@latest

.PHONY: get-install-tailwindcss
get-install-tailwindcss: ## Installs the tailwindcss cli
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
	chmod +x tailwindcss-linux-x64
	mv tailwindcss-linux-x64 tailwindcss

# --- Development Tasks ---
.PHONY: tailwind-watch
tailwind-watch: ## compile tailwindcss and watch for changes
	./tailwindcss -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT) --watch

.PHONY: tailwind-build
tailwind-build: ## one-time compile tailwindcss styles
	./tailwindcss -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT)

.PHONY: templ-generate
templ-generate: ## generate go code from templ files
	templ generate

.PHONY: templ-watch
templ-watch: ## watch templ files and generate on change
	templ generate --watch

# --- Build & Run ---
.PHONY: build
build: tailwind-build templ-generate ## compile assets and build the main application
	@echo "Building $(APP_NAME) binary..."
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_PATH) $(MAIN_PKG)
	@echo "Build complete: $(BINARY_PATH)"

.PHONY: run
run: build ## build and run the main application
	@echo "Running $(APP_NAME)..."
	./$(BINARY_PATH)

.PHONY: watch
watch: ## build and watch the project with air
	@echo "Watching for changes with air..."
	air

# --- Cleanup ---
.PHONY: clean
clean: ## remove build artifacts
	@echo "Cleaning up..."
	rm -f $(BINARY_PATH)