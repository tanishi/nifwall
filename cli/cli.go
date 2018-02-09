package cli

import (
	"flag"
	"fmt"
	"io"
)

const (
	exitCodeOK = iota
	exitCodeParseFlagError
)

// CLI is structure for cli tool
type CLI struct {
	OutStream, ErrStream io.Writer
}

const version string = "v0.1.0"

// Run execute cli
func (c *CLI) Run(args []string) int {
	var v bool
	flags := flag.NewFlagSet("nifwall", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)
	flags.BoolVar(&v, "version", false, "Print version information and quit")

	if err := flags.Parse(args[1:]); err != nil {
		return exitCodeParseFlagError
	}

	if v {
		fmt.Fprintf(c.ErrStream, "nifwall version %s\n", version)
		return exitCodeOK
	}

	return exitCodeOK
}
