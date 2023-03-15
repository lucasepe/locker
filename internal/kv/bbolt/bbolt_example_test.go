package bbolt

import (
	"fmt"
	"log"
	"os"

	"github.com/lucasepe/locker/internal/kv"
)

func ExampleStore_PutOne() {
	opt := Options{
		Path:  tempfile(),
		Codec: kv.NewCryptoCodec("HELLO!"),
	}

	sto, err := NewStore(opt)
	if err != nil {
		panic(err)
	}
	defer os.Remove(opt.Path)
	defer sto.Close()

	namespace, key := "google.com", "password"
	err = sto.PutOne(namespace, key, "abbracadabbra")
	if err != nil {
		panic(err)
	}

	got, err := sto.GetOne(namespace, key)
	if err != nil {
		panic(err)
	}
	fmt.Println(got)

	// Output:
	// abbracadabbra
}

// tempfile returns a temporary file path.
func tempfile() string {
	f, err := os.CreateTemp("", "bolt-")
	if err != nil {
		log.Fatalf("Could not create temp file: %s", err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		log.Fatal(err)
	}
	return f.Name()
}
