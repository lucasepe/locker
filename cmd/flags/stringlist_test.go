package flags

import (
	"flag"
	"testing"
)

func TestStringList(t *testing.T) {
	fv := StringList{}

	var fs flag.FlagSet
	fs.Var(&fv, "key", "")

	err := fs.Parse([]string{"-key", "url", "-key", "username", "-key", "password"})
	if err != nil {
		t.Fail()
	}

	want := "[url username password]"
	got := fv.String()
	if got != want {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}
