package cmd_test

import (
	"bytes"
	"testing"

	"github.com/paxthemax/resolver"
	"github.com/paxthemax/resolver/cmd/resolver-cli/cmd"
)

func TestVersionCmd(t *testing.T) {
	var outputBuffer bytes.Buffer
	if err := newCommand(t,
		cmd.WithArgs("version"),
		cmd.WithOutput(&outputBuffer),
	).Execute(); err != nil {
		t.Fatal(err)
	}

	want := resolver.Version + "\n"
	got := outputBuffer.String()
	if got != want {
		t.Errorf("Got version output %q, want %q", got, want)
	}
}
