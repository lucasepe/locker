package flags

import (
	"flag"
	"testing"
)

func TestLabelFlag(t *testing.T) {
	fv := LabelFlag{}

	var fs flag.FlagSet
	fs.Var(&fv, "key", "")

	err := fs.Parse([]string{"-key", "è un mio segreto!"})
	if err != nil {
		t.Fail()
	}

	want := "unMioSegreto"
	got := fv.String()
	if got != want {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}
