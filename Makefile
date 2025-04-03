.PHONY: build
VERSION=$$(git describe --tags --abbrev=0)
SOURCE_FILE := main.go
BUILD_DIR := build
RELEASE_DIR := release
OS_LIST := linux darwin windows
ARCH_LIST := amd64 386 arm arm64
INSTALL_FILE := /usr/local/bin/executor
EXE_NAME := executor
STYLE_LIST := line blink arrow star circle square outline bar o cursor dots

version:
	@echo $(VERSION)

test:
	go test ./...
test-v:
	go test -v ./...

run:
	@go run main.go run --desc "Short run test" -c "sleep 2; echo \"!Hola, Mundo!\""; \

run-which:
	@go run main.go which -c ls
	@echo "Running silently OK"
	@go run main.go which --silent -c ls
	@echo "Running silently FAIL"
	@go run main.go which --silent -c non-existing-command

demo:
	@for style in $(STYLE_LIST); do \
		DESC="Making something with '$$style' style"; \
		SECS=$$(($${RANDOM} % 3 + 1)); \
		echo "\033[90m#> executor run --desc \"$${DESC}\" -st $$style -c \"sleep $${SECS}; echo Hello;\"\033[0m" | pv -qL 60; \
		go run main.go run --desc "$${DESC}" -st $$style -c "sleep $${SECS}; echo Hello"; \
	done
	@echo "\033[90mexecutor run --desc \"Not interactive and no color test\" --nc -st bar -c \"sleep 1; echo Hello\"\033[0m" | pv -qL 60
	@go run main.go run --desc "Not interactive and no color test" --nc -st bar -c "sleep 1; echo Hello"
	@echo "\033[90m#> executor which -st $$style -c \"ls\"\033[0m" | pv -qL 60;
	@go run main.go which -c "ls"

run-long:
	@go run main.go run --desc "Long run test" -c "sleep 63; echo Hello";
run-short:
	@go run main.go run --desc "Short run test" -c "sleep 0.2; echo Hello";
run-with-slow-output:
	@go run main.go run --desc "Run with output" -c "cat Makefile | pv -qL 150"
run-help:
	@go run main.go run --help

run-show-env:
	@go run main.go run -se --desc "Show env" -c "echo '¡Hola, Mundo!'"

run-output:
	@echo
	@go run main.go run --desc "No error, show output" -so -c "sleep 1; echo '¡Hola, Mundo!'; echo 'No errors' >&2"
	@echo
	@go run main.go run --desc "Error, show output" -so -c "sleep 1; echo '¡Hola, Mundo!'; echo 'This is an error' >&2; exit 1"

build: clean lint test
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(EXE_NAME) main.go
	@echo "Built to $(BUILD_DIR)/$(EXE_NAME)"

clean:
	@rm -rf $(BUILD_DIR)
	@rm -rf $(RELEASE_DIR)

lint:
	@echo "Linting..."
	@go tool staticcheck ./...
	@golangci-lint run ./...

install: build
	@cp $(BUILD_DIR)/$(EXE_NAME) $(INSTALL_FILE)
	@echo "Installed to $(INSTALL_FILE)"

uninstall:
	@rm -rf $(INSTALL_FILE)
	@echo "Removed $(INSTALL_FILE)"

release: clean
	@mkdir -p $(RELEASE_DIR)
	@for os in $(OS_LIST); do \
		for arch in $(ARCH_LIST); do \
			f="$(EXE_NAME)"; \
			crunched_name="$(EXE_NAME)_$${os}-$${arch}_$(subst .,_,${VERSION}).tar.bz"; \
			if [ "$$os" = "windows" ]; then \
				f="$$f.exe"; \
				crunched_name="$(EXE_NAME)_$${os}-$${arch}_$(subst .,_,${VERSION}).zip"; \
			fi; \
			if [[ "$$os/$$arch" != "darwin/arm" && "$$os/$$arch" != "darwin/386" ]]; then \
				echo "Building $$os/$$arch --> $$f"; \
				GOOS=$$os GOARCH=$$arch go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(RELEASE_DIR)/$$f $(SOURCE_FILE); \
				echo "Compressing $$f --> $$crunched_name"; \
				if [ "$$os" = "windows" ]; then \
					pushd $(RELEASE_DIR); \
					zip -r $$crunched_name $$f; \
					popd; \
				else \
					tar -C $(RELEASE_DIR) -cjf $(RELEASE_DIR)/$$crunched_name $$f; \
				fi; \
				rm $(RELEASE_DIR)/$$f; \
			fi; \
		done; \
	done
	@for f in $(RELEASE_DIR)/*; do \
		i=$$(basename $$f); \
		size=$$(du -sh $$f | cut -f 1); \
		echo "$$size\t$$(sha256sum $$f)"; \
	done
