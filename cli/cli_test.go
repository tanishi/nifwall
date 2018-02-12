package cli

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
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

	expected := fmt.Sprintf("nifwall version %s\n", version)
	actual := errStream.String()

	if expected != actual {
		t.Errorf("expected: %v, but: %v", expected, actual)
	}

	outStream.Reset()
	errStream.Reset()
	args = strings.Split("nifwall list", " ")

	status = cli.Run(args)
	if status != exitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
	}

	expected = fmt.Sprintf("list\n")
	actual = outStream.String()

	if expected != actual {
		t.Errorf("expected: %v, but: %v", expected, actual)
	}

	outStream.Reset()
	errStream.Reset()
	args = strings.Split("nifwall update", " ")

	status = cli.Run(args)
	if status != exitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
	}

	expected = fmt.Sprintf("update\n")
	actual = errStream.String()

	if expected != actual {
		t.Errorf("expected: %v, but: %v", expected, actual)
	}

	outStream.Reset()
	errStream.Reset()
	args = strings.Split("nifwall apply", " ")

	status = cli.Run(args)
	if status != exitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
	}

	expected = fmt.Sprintf("apply\n")
	actual = outStream.String()

	if expected != actual {
		t.Errorf("expected: %v, but: %v", expected, actual)
	}
}
