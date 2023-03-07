package flags

import (
	"flag"
	"testing"
)

func TestBoxFlag(t *testing.T) {
	fv := BoxFlag{}

	var fs flag.FlagSet
	fs.Var(&fv, "box", "")

	err := fs.Parse([]string{"-box", "my-secrets"})
	if err != nil {
		t.Fail()
	}

	want := "my-secrets"
	got := fv.String()
	if got != want {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}
