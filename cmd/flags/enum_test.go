package flags

import (
	"flag"
	"reflect"
	"testing"
)

func TestEnum(t *testing.T) {
	fv := Enum{Choices: []string{
		"uri",
		"username",
		"password",
		"totp",
		"token",
		"query_string",
	}}
	var fs flag.FlagSet
	fs.Var(&fv, "type", "")

	err := fs.Parse([]string{"-type", "totp", "-type", "query_string"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, "query_string") {
		t.Fail()
	}
}
