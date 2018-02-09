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
		OutStream: outStream,
		ErrStream: errStream,
	}
	args := strings.Split("nifwall -version", " ")

	status := cli.Run(args)
	if status != exitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
	}

	expected := fmt.Sprintf("nifwall version %s", version)
	actual := errStream.String()

	if expected == actual {
		t.Errorf("expected: %v, but: %v", expected, actual)
	}
}
