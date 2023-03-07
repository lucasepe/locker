// This Go package removes ANSI escape codes from strings.
//
// Ideally, we would prevent these from appearing in any text we want to process.
// However, sometimes this can't be helped, and we need to be able to deal with that noise.
// This will use a regexp to remove those unwanted escape codes.
//
// Credits to: https://github.com/acarl005/stripansi
package stripansi

import (
	"regexp"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func Strip(str string) string {
	return re.ReplaceAllString(str, "")
}
