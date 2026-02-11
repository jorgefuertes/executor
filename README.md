# ![Executor Logo](assets/executor.svg)

[![Go Report Card](https://goreportcard.com/badge/github.com/jorgefuertes/executor)](https://goreportcard.com/report/github.com/jorgefuertes/executor)
[![GoDoc](https://godoc.org/github.com/jorgefuertes/executor?status.svg)](https://godoc.org/github.com/jorgefuertes/executor)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/jorgefuertes/executor/total)
![GitHub Release](https://img.shields.io/github/v/release/jorgefuertes/executor)
![License](https://img.shields.io/github/license/jorgefuertes/executor)

## Execute commands in a fancy way

A command-line tool to run shell commands with visual feedback, timing, and clean output. Perfect for build scripts, CI/CD pipelines, and any workflow where you need clear visual confirmation of command execution.

**Not a replacement** for `make` or `task` - it's a complement. When your build output is too verbose, `executor` helps you see at a glance if everything went well.

## Features

✨ **Visual Feedback** - Spinners and colorized output for better UX  
🎯 **Smart Output** - Hides log pollution, shows only what matters  
⏱️ **Timing** - Displays execution time for each command  
🔍 **Multiple Modes** - Run commands, check ports, validate web endpoints, verify executables  
🎨 **Auto-detection** - Switches to plain output when not in a TTY  
🚀 **Fast** - Written in Go, single binary with no dependencies  

### Output Shading

Executor hides standard and error output by default, displaying it only on error or when explicitly requested via flags (`--show-output`, `-so`).

### Non-Interactive Mode

Automatic terminal detection - uses plain output (no colors/spinners) when output is redirected or when `--no-interactive` flag is set.

## Demo

![Demo](./assets/demo.gif)

## Installation

### macOS (Homebrew)

**Recommended:**

~~~bash
brew install jorgefuertes/executor/executor
~~~

Or tap first:

~~~bash
brew tap jorgefuertes/executor
brew install executor
~~~

### Linux and Others

Using curl:

~~~bash
curl -Lo - https://raw.githubusercontent.com/jorgefuertes/executor/refs/heads/main/scripts/install.sh | sh
~~~

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/jorgefuertes/executor/releases).

### As a Go Tool

Add to your `go.mod`:

~~~go
tool (
    github.com/jorgefuertes/executor
)
~~~

Then:

~~~bash
go mod tidy
go tool executor --help
~~~

## Usage

### Getting Help

~~~bash
executor --help
executor run --help  # Help for specific command
~~~

### Run Commands

Basic command execution:

~~~bash
executor run --desc "Build project" -c "npm run build"
~~~

Show output on success:

~~~bash
executor run --show-output --desc "Run tests" -c "go test ./..."
~~~

Using short flags:

~~~bash
executor run -so -d "Deploy" -c "./deploy.sh"
~~~

### Check Executable Existence

~~~bash
executor which --desc "Check Docker" -c docker
executor which -d "Check Node.js" -c node
~~~

### Check Port Availability

~~~bash
executor port --desc "Check PostgreSQL" -p 5432
executor port -d "Check web server" -p 8080 --host localhost
~~~

### Validate Web Endpoints

~~~bash
executor web --desc "Check API" -u https://api.example.com/health
executor web -d "Verify homepage" -u https://example.com --status 200
~~~

## Available Commands

| Command | Description | Example |
|---------|-------------|---------|
| `run` | Execute a shell command or script | `executor run -c "make build"` |
| `which` | Verify executable exists in PATH | `executor which -c docker` |
| `port` | Check if a port is open | `executor port -p 5432` |
| `web` | Validate HTTP endpoint responds | `executor web -u https://api.com` |

### Common Flags

- `--desc, -d` - Description to display
- `--show-output, -so` - Show command output even on success
- `--no-color, --nc` - Disable colors
- `--no-interactive, --ni` - Disable spinners (for CI/CD)
- `--help, -h` - Show help

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

GPL-3.0-or-later - See [LICENSE](LICENSE) file for details.

## Authors

- **Idea**: Marcos Gómez
- **Main Developer**: [Jorge Fuertes (Queru)](https://github.com/jorgefuertes)

---

⭐ If you find this tool useful, please consider giving it a star on GitHub!
