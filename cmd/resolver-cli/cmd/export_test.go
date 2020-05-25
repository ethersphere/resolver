package cmd

import (
	"io"
)

type (
	Command = command
	Option  = option
)

var (
	NewCommand = newCommand

	// Avoid unused linter errors until functions are used:
	_ = WithConfgFile
	_ = WithInput
	_ = WithErrorOutput
)

func WithConfgFile(f string) func(c *Command) {
	return func(c *Command) {
		c.configFile = f
	}
}

func WithArgs(a ...string) func(c *Command) {
	return func(c *Command) {
		c.root.SetArgs(a)
	}
}

func WithInput(r io.Reader) func(c *Command) {
	return func(c *Command) {
		c.root.SetIn(r)
	}
}

func WithOutput(w io.Writer) func(c *Command) {
	return func(c *Command) {
		c.root.SetOut(w)
	}
}

func WithErrorOutput(w io.Writer) func(c *Command) {
	return func(c *Command) {
		c.root.SetErr(w)
	}
}
