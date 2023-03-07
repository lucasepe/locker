package cmd

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/lucasepe/locker/cmd/app"
)

func TestCmdImport(t *testing.T) {
	defer os.Remove(testArchivePath())

	os.Setenv(app.EnvSecret, testSecret)

	out := bytes.NewBufferString("")
	err := runCmdImport(out)
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(out.String())
	want := "successfully imported 3 documents"
	if got != want {
		t.Fatalf("expected:%v, got:%v", want, got)
	}
}

func runCmdImport(output io.Writer) error {
	op := newCmdImport()

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(output)

	op.SetFlags(fs)

	args := []string{
		"-n", testStore,
		"-f", "../testdata/sample.yaml",
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	return op.Execute(fs)
}
