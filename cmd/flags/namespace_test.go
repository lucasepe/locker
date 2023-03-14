package flags

import (
	"flag"
	"testing"
)

func TestNamespaceFlag(t *testing.T) {
	fv := NamespaceFlag{}

	var fs flag.FlagSet
	fs.Var(&fv, "namespace", "")

	err := fs.Parse([]string{"-namespace", "Google.com"})
	if err != nil {
		t.Fail()
	}

	want := "google-com"
	got := fv.String()
	if got != want {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}
