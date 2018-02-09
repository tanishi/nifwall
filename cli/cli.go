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
	OutStream io.Writer
	ErrStream io.Writer
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

	if len(args) < 2 {
		fmt.Fprint(c.ErrStream, "list or update or apply")
		return exitCodeOK
	}

	if v {
		fmt.Fprintf(c.ErrStream, "nifwall version %s\n", version)
		return exitCodeOK
	}

	switch args[1] {
	case "list":
		fmt.Fprintf(c.OutStream, "list")
	case "update":
		fmt.Fprintf(c.ErrStream, "update")
	case "apply":
		fmt.Fprintf(c.OutStream, "apply")
	default:
		flag.PrintDefaults()
	}

	return exitCodeOK
}
