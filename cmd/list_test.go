package cmd

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCmdList(t *testing.T) {
	defer os.Remove(testArchivePath())

	os.Setenv(EnvSecret, testSecret)

	out := bytes.NewBufferString("")
	if err := runCmdPut(out, "user name", "Pino Latino"); err != nil {
		t.Fatal(err)
	}
	if err := runCmdPut(out, "password", "Non te la dico"); err != nil {
		t.Fatal(err)
	}

	out.Reset()
	if err := runCmdList(out); err != nil {
		t.Fatal(err)
	}

	want := []string{
		"password",
		"user_name",
	}

	got := strings.Fields(strings.TrimSpace(out.String()))
	for i, el := range got {
		got[i] = strings.TrimSpace(el)
	}

	if !cmp.Equal(want, got) {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}

func runCmdList(output io.Writer) error {
	op := newCmdList()

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(output)

	op.SetFlags(fs)

	err := fs.Parse([]string{
		"-n", testNamespace,
		"-s", testStore,
	})
	if err != nil {
		return err
	}

	return op.Execute(fs)
}
