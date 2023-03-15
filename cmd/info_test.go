package cmd

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCmdInfo(t *testing.T) {
	defer os.Remove(testArchivePath())

	out := bytes.NewBufferString("")
	err := runCmdInfo(out, "1.0.0", "8888")
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(out.String())
	want := "Locker 1.0.0 (build: 8888) "
	if !strings.HasPrefix(got, want) {
		t.Fatalf("expected prefix: %v, got: %v", want, got)
	}
}

func runCmdInfo(output io.Writer, ver, bld string) error {
	op := newCmdInfo(ver, bld)

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(output)

	op.SetFlags(fs)

	err := fs.Parse([]string{})
	if err != nil {
		return err
	}

	return op.Execute(fs)
}
