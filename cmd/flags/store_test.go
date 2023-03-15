package flags

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/lucasepe/xdg"
)

func TestStore(t *testing.T) {
	fv := Store{
		BaseDir: filepath.Join(xdg.ConfigDir(), "Locker"),
	}

	var fs flag.FlagSet
	fs.Var(&fv, "store", "")

	err := fs.Parse([]string{"-store", "test"})
	if err != nil {
		t.Fail()
	}

	want := filepath.Join(xdg.ConfigDir(), "Locker", "test.db")
	got := fv.String()
	if got != want {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}
