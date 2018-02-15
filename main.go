package main

import (
	"os"

	nifcloud "github.com/tanishi/go-nifcloud"
	"github.com/tanishi/nifwall/cli"
	nifwall "github.com/tanishi/nifwall/lib"
)

func main() {
	endpoint := os.Getenv("NIFCLOUD_ENDPOINT")
	accessKey := os.Getenv("NIFCLOUD_ACCESSKEY")
	secretAccessKey := os.Getenv("NIFCLOUD_SECRET_ACCESSKEY")

	nifclient, _ := nifcloud.NewClient(endpoint, accessKey, secretAccessKey)
	client := nifwall.NewClient(nifclient)

	c := &cli.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
		Client:    client,
	}

	os.Exit(c.Run(os.Args))
}
