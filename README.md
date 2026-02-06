# Executor

| ![Executor Logo](assets/executor.svg) |
|:--:|

Run commands in a fancy way.

## About executor

This is a small program just made to run any shell commands with visual feedback, timing and some other features, giving it a better and professional look.

This is not intended to be an alternative to `make` or `task`, this is just a complement. When your `make` output is too much to understand, `executor` helps you to see, in a single look, if everything goes ok.

### Output shading

Executor hides the standard and error output, displaying it only on error or as instructed via command line flags.

### Not interactive mode

Automatic detection of terminal type, acts as a simple log without colors or spinners when the output is not a terminal or the flag `--no-color` is activated.

## Demo

![Demo](./assets/demo.gif)

## Installation

### MacOS-X (homebrew)

~~~bash
brew tap jorgefuertes/executor
brew install executor
~~~

### Linux and others (curl)

~~~bash
curl -Lo - https://raw.githubusercontent.com/jorgefuertes/executor/refs/heads/main/scripts/install.sh | sh
~~~

## Quick usage

### Getting help

~~~bash
executor --help
~~~

### Command run

~~~bash
executor run --desc "Just a greeting" -c "sleep 3; echo 'Hello World!'"
~~~

### Showing the output

~~~bash
executor run -so --desc "Just a greeting" -c "sleep 3; echo 'Hello World!'"
~~~

## Run modes

You can use executor, at the moment, to:

- Run a command or a script (`run`).
- Check for the existence of an executable (`which`).
- Check if a port is open (`port`).
- Check if a web page is responding successfully (`web`).

Check the included help `--help` for more info.

## Authors

- Idea: Marcos Gómez.
- Main developer: [Jorge Fuertes AKA Queru](https://github.com/jorgefuertes).
