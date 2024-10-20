package commands

import "errors"

var ErrCommandNotFound = errors.New("command not found")
var ErrEmptyCommand = errors.New("empty command")
var ErrExecutionFailed = errors.New("execution failed")
