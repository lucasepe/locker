package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lucasepe/locker/cmd/app"
)

const (
	testStore  = "test"
	testBox    = "stuffs"
	testSecret = "Abbracadabbra"
)

func TestCmdAdd(t *testing.T) {
	defer os.Remove(testArchivePath())

	os.Setenv(app.EnvSecret, testSecret)

	out := bytes.NewBufferString("")
	if err := runCmdPut(out, "user", "Pino Latino"); err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(out.String())
	want := "item successfully stored"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("expected prefix: %v, got: %v", want, got)
	}
}

func runCmdPut(output io.Writer, k, v string) error {
	op := newCmdPut()

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(output)

	op.SetFlags(fs)

	err := fs.Parse([]string{
		"-b", testBox,
		"-n", testStore,
		"-l", k,
		v,
	})
	if err != nil {
		return err
	}

	return op.Execute(fs)
}

func testArchivePath() string {
	return filepath.Join(app.Dir(), fmt.Sprintf("%s.db", testStore))
}
