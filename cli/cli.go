package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	nifcloud "github.com/tanishi/go-nifcloud"
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

	endpoint := os.Getenv("NIFCLOUD_ENDPOINT")
	accessKey := os.Getenv("NIFCLOUD_ACCESSKEY")
	secretAccessKey := os.Getenv("NIFCLOUD_SECRET_ACCESSKEY")

	nifwall.Client, _ = nifcloud.NewClient(endpoint, accessKey, secretAccessKey)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch args[1] {
	case "list":

		instances, err := nifwall.ListInappropriateInstances(ctx, "fwName")

		if err != nil {
			fmt.Fprint(c.ErrStream, err)
			return exitCodeError
		}

		fmt.Fprint(c.OutStream, instances)

	case "update":
		fmt.Fprintf(c.ErrStream, "update\n")
	case "apply":
		fmt.Fprintf(c.OutStream, "apply\n")
	default:
		flag.PrintDefaults()
	}

	return exitCodeOK
}
