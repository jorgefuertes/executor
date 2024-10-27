package commands

import "errors"

var (
	ErrCommandNotFound = errors.New("command not found")
	ErrEmptyCommand    = errors.New("empty command")
	ErrExecutionFailed = errors.New("execution failed")
	ErrEnvNotFound     = errors.New("env not found")
)
