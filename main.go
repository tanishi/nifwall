package main

import (
	"os"

	"github.com/tanishi/nifwall/cli"
)

func main() {
	c := &cli.CLI{OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(c.Run(os.Args))
}
