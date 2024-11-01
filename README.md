# Executor

Run commands in a fancy way.

## About executor

This is a small program just made to run any shell commands with visual feedback, timing and some other features, giving it a better and professional look.

### Output shading

Executor hides the standard and error output, showing it on error or as instructed via command line flags.

### Not interactive mode

Automatic detection of terminal type, acts as a simple log without colors or spinners when the output is not a terminal or the flag `--no-color` is activated.

## Demo

![Demo](./assets/demo.gif)

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
## English usage check

Upon checking this document's English, there are a few corrections to note:

- "This is an small program" should be "This is a small program"
- "giving a better and pro looking" should be "giving it a better and professional look"
- "showing it on error or as instructed" should be "showing output on error or as instructed"
- "acts as a simple log" could be "functions as a simple log"
- "reponding successfuly" should be "responding successfully"

Overall the document is quite readable with only minor grammatical improvements needed.
