# Executor - Copilot Instructions

## About

Executor is a Go CLI tool that runs shell commands with visual feedback (spinners, timing, colored output). It's designed to make command execution output cleaner and more professional, particularly useful for complex build/test pipelines.

## Build, Test, and Lint

### Build
```bash
make build       # Builds to build/executor (runs lint + test first)
make install     # Builds and installs to /usr/local/bin/executor
```

### Testing
```bash
make test        # Run all tests
make test-v      # Run tests with verbose output
go test ./...    # Standard Go test invocation
```

Run individual test files:
```bash
go test -v ./internal/commands -run TestWhich
```

### Linting
```bash
make lint        # Runs staticcheck and golangci-lint
```

The project uses:
- `go tool staticcheck` (via go.mod tools directive)
- `golangci-lint`

## Architecture

### Package Structure

- **main.go**: CLI entry point using `urfave/cli/v2`. Defines 4 commands: `run`, `which`, `port`, `web`
- **internal/config**: Parses CLI flags into a `Config` struct passed to all commands
- **internal/commands**: Command implementations (`run.go`, `which.go`, `port.go`, `web.go`)
- **internal/terminal**: Terminal abstraction handling colors, spinners, progress display

### Key Flow

1. main.go defines CLI structure with urfave/cli/v2
2. Each command action wrapped by `newActionFunc()` which:
   - Creates `Config` from CLI context
   - Optionally prints config if `--show-config` flag set
   - Calls command function with config
3. Command functions instantiate `terminal.Term` and use its methods for output
4. Terminal package handles:
   - Color/interactive mode detection (auto-disables if not TTY)
   - Spinner animations (11 styles available)
   - Progress tracking with buffers for stdout/stderr
   - Cleanup (restoring cursor visibility)

### Environment File Loading

The `run` command supports `.env` file loading with recursive directory search:
- Searches up to 5 levels (configurable via `--env-recurse-levels`)
- Uses `github.com/joho/godotenv` package
- Env vars merged into command environment
- Can be disabled with `--env-file none`

### Terminal Abstraction

The `terminal.Term` type auto-detects terminal capabilities:
- Checks if output is a TTY
- Disables colors/spinners in non-interactive environments
- Can be forced via `--no-color` and `--no-interactive` flags
- Progress type buffers stdout/stderr, displays on error by default

## Key Conventions

### Config Pattern

All commands receive a `*config.Config` struct. Never parse CLI flags directly in command functions—always access via config fields.

### Terminal Output

- Use `terminal.Term` methods for all output, not direct `fmt.Print`
- Always call `defer t.CleanUp()` after creating Terminal to restore cursor
- Commands should return errors; terminal formatting is separate concern

### Spinner Styles

Available styles defined in `internal/terminal/spinner.go`:
- dots (default), arrow, star, circle, square, outline, line, bar, o, cursor, blink

### Version Injection

Version is set via ldflags during build:
```bash
-ldflags="-s -w -X main.version=$(VERSION)"
```

The version comes from git tags (`git describe --tags --abbrev=0`).

### Cross-compilation

The Makefile's `release` target builds for:
- OS: linux, darwin, windows
- Arch: amd64, 386, arm, arm64
- Excludes: darwin/arm, darwin/386

Output is compressed to `.tar.bz` (Unix) or `.zip` (Windows).

## Testing

Only 2 test files exist:
- `internal/commands/env_test.go`
- `internal/commands/which_test.go`

Use `testify/assert` for assertions.
