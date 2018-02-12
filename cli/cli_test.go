package cli

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	cases := []struct {
		args     string
		expected string
	}{
		{"nifwall -version", fmt.Sprintf("nifwall version %s\n", version)},
		{"nifwall list", fmt.Sprintf("list\n")},
		{"nifwall update", fmt.Sprintf("update\n")},
		{"nifwall apply", fmt.Sprintf("apply\n")},
		{"nifwall", fmt.Sprintf("list or update or apply")},
	}

	for _, c := range cases {
		outStream := new(bytes.Buffer)

		cli := &CLI{
			OutStream: outStream,
			ErrStream: outStream,
		}

		status := cli.Run(strings.Split(c.args, " "))

		if status != exitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
		}

		actual := outStream.String()

		if c.expected != actual {
			t.Errorf("expected: %v, but: %v", c.expected, actual)
		}
	}
}
