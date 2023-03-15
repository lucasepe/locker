package cmd

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCmdImport(t *testing.T) {
	defer os.Remove(testArchivePath())

	os.Setenv(EnvSecret, testSecret)

	out := bytes.NewBufferString("")
	err := runCmdImport(out)
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(out.String())
	want := "successfully imported 2 documents"
	if got != want {
		t.Fatalf("expected: %s, got: %s", want, got)
	}
}

func runCmdImport(output io.Writer) error {
	op := newCmdImport()

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(output)

	op.SetFlags(fs)

	args := []string{
		"-s", testStore,
		"-f", "../testdata/sample.yaml",
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	return op.Execute(fs)
}
