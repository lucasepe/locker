package kv

import (
	"fmt"
	"os"

	"github.com/lucasepe/locker/cmd/app"
)

func ExampleCodec_Marshal() {
	os.Setenv(app.EnvSecret, "MAGIK")

	codec := NewCryptoCodec()
	enc, err := codec.Marshal([]byte("Hello World!"))
	if err != nil {
		panic(err)
	}

	dec, err := codec.Unmarshal(enc)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(dec))

	// Output:
	// Hello World!
}
