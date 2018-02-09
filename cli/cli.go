package cli

import (
	"flag"
	"fmt"
	"io"
)

const (
	ExitCodeOK = iota
	ExitCodeParseFlagError
)

type CLI struct {
	OutStream, ErrStream io.Writer
}

const Version string = "v0.1.0"

func (c *CLI) Run(args []string) int {
	var version bool
	flags := flag.NewFlagSet("nifwall", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)
	flags.BoolVar(&version, "version", false, "Print version information and quit")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	if version {
		fmt.Fprintf(c.ErrStream, "nifwall version %s\n", Version)
		return ExitCodeOK
	}

	return ExitCodeOK
}
