package cli

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRun_versionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: outStream,
		errStream: errStream,
	}
	args := strings.Split("nifwall -version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("nifwall version %s", Version)
	actual := errStream.String()

	if expected == actual {
		t.Errorf("expected: %v, but: %v", expected, actual)
	}
}
