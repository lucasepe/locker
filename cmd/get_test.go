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

func TestCmdGetOne(t *testing.T) {
	defer os.Remove(testArchivePath())

	os.Setenv(app.EnvSecret, testSecret)

	out := bytes.NewBufferString("")
	err := runCmdPut(out, "user name", "pinco.pallo@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	out.Reset()
	err = runCmdGet(out, "userName")
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(out.String())
	want := "pinco.pallo@gmail.com"
	if got != want {
		t.Fatalf("expected: %s, got: %s", want, got)
	}
}

func TestCmdGetAll(t *testing.T) {
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
	if err := runCmdGet(out, ""); err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(out.String())
	want := `password: magick
userName: pinco.pallo@gmail.com`
	if got != want {
		t.Fatalf("expected:%v, got:%v", want, got)
	}
}

func runCmdGet(output io.Writer, key string) error {
	op := newCmdGet()

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
