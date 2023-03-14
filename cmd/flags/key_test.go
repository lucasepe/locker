package flags

import (
	"flag"
	"testing"
)

func TestKeyFlag(t *testing.T) {
	fv := KeyFlag{}

	var fs flag.FlagSet
	fs.Var(&fv, "key", "")

	err := fs.Parse([]string{"-key", "è un mio segreto!"})
	if err != nil {
		t.Fail()
	}

	want := "un_mio_segreto"
	got := fv.String()
	if got != want {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}
