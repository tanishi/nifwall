package cli

import (
	"context"
	"flag"
	"fmt"
	"io"

	nifwall "github.com/tanishi/nifwall/lib"
)

const (
	exitCodeOK = iota
	exitCodeParseFlagError
	exitCodeError
)

// CLI is structure for cli tool
type CLI struct {
	OutStream io.Writer
	ErrStream io.Writer
	Client    *nifwall.Client
}

type firewallList []string

func (f *firewallList) String() string {
	return "firewall list"
}

func (f *firewallList) Set(s string) error {
	*f = append(*f, s)
	return nil
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch args[1] {
	case "list":
		listFlags := flag.NewFlagSet("list", flag.ContinueOnError)

		var fws firewallList
		listFlags.Var(&fws, "fw", "specify firewall")

		if err := listFlags.Parse(args[2:]); err != nil {
			return exitCodeParseFlagError
		}

		instances, err := c.Client.ListInappropriateInstances(ctx, fws)

		if err != nil {
			fmt.Fprint(c.ErrStream, err)
			return exitCodeError
		}

		fmt.Fprint(c.OutStream, instances)

	case "update":
		updateFlags := flag.NewFlagSet("update", flag.ContinueOnError)

		var f string
		updateFlags.StringVar(&f, "f", "nifwall", "specify firewall")

		if err := updateFlags.Parse(args[2:]); err != nil {
			return exitCodeParseFlagError
		}

		err := c.Client.UpdateFirewall(ctx, f)

		if err != nil {
			fmt.Fprint(c.ErrStream, err)
			return exitCodeError
		}

	case "apply":
		applyFlags := flag.NewFlagSet("apply", flag.ContinueOnError)

		var fw string
		applyFlags.StringVar(&fw, "fw", "nifwall", "specify firewall")

		if err := applyFlags.Parse(args[2:]); err != nil {
			return exitCodeParseFlagError
		}

		c.Client.RegisterInstancesWithSecurityGroup(ctx, fw, applyFlags.Args())

	default:
		flag.PrintDefaults()
	}

	return exitCodeOK
}
