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

func TestCmdDeleteOne(t *testing.T) {
	defer os.Remove(testArchivePath())

	os.Setenv(app.EnvSecret, testSecret)

	out := bytes.NewBufferString("")
	err := runCmdPut(out, "user name", "pinco.pallo@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	out.Reset()
	err = runCmdDelete(out, "userName")
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(out.String())
	want := "item successfully deleted"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("expected prefix: %v, got: %v", want, got)
	}
}

func TestCmdDeleteAll(t *testing.T) {
	defer os.Remove(testArchivePath())

	os.Setenv(app.EnvSecret, testSecret)

	out := bytes.NewBufferString("")
	err := runCmdPut(out, "user name", "pinco.pallo@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	err = runCmdPut(out, "password", "magick")
	if err != nil {
		t.Fatal(err)
	}

	out.Reset()
	if err := runCmdDelete(out, ""); err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(out.String())
	want := "box successfully deleted"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("expected prefix: %v, got: %v", want, got)
	}
}

func runCmdDelete(output io.Writer, key string) error {
	op := newCmdDelete()

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(output)

	op.SetFlags(fs)

	args := []string{
		"-b", testBox,
		"-n", testStore,
	}
	if len(key) > 0 {
		args = append(args, "-l", key)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	return op.Execute(fs)
}
