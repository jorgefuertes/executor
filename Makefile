.PHONY: build
VERSION=v1.0.0
SOURCE_FILE := main.go
BUILD_DIR := build
RELEASE_DIR := release
OS_LIST := linux darwin windows
ARCH_LIST := amd64 386 arm arm64
INSTALL_FILE := ~/bin/executor
EXE_NAME := executor

version:
	@echo $(VERSION)

run:
	go run main.go  -desc test -show-env -show-output stdout -show-on-err both echo "!Hola, Mundo!"

build: clean lint
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w" -o $(BUILD_DIR)/$(EXE_NAME) main.go
	@echo "Built to $(BUILD_DIR)/$(EXE_NAME)"

clean:
	@rm -rf $(BUILD_DIR)
	@rm -rf $(RELEASE_DIR)

lint:
	@echo "Linting..."
	@staticcheck ./...
	@golangci-lint run ./...

install: build
	@cp $(BUILD_DIR)/$(EXE_NAME) $(INSTALL_FILE)
	@echo "Installed to $(INSTALL_FILE)"

uninstall:
	@rm -rf $(INSTALL_FILE)
	@echo "Removed $(INSTALL_FILE)"

release: clean lint
	@mkdir -p $(RELEASE_DIR)
	@for os in $(OS_LIST); do \
		for arch in $(ARCH_LIST); do \
			f="$(EXE_NAME)"; \
			tar_name="$(EXE_NAME)_$${os}-$${arch}_$(subst .,_,${VERSION}).tar.bz"; \
			if [ "$$os" = "windows" ]; then \
				f="$$f.exe"; \
			fi; \
			if [[ "$$os/$$arch" != "darwin/arm" && "$$os/$$arch" != "darwin/386" ]]; then \
				echo "Building $$os/$$arch --> $$f"; \
				GOOS=$$os GOARCH=$$arch go build -o $(RELEASE_DIR)/$$f $(SOURCE_FILE); \
				echo "Compressing $$f --> $$tar_name"; \
				tar -C $(RELEASE_DIR) -cjf $(RELEASE_DIR)/$$tar_name $$f; \
				rm $(RELEASE_DIR)/$$f; \
			fi; \
		done; \
	done
	@ls -sSFhC1 release