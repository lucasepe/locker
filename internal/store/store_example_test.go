package store_test

import (
	"fmt"
	"os"
	"time"

	"github.com/lucasepe/locker/internal/store"
	"go.etcd.io/bbolt"
)

// Show we can put an item in a bucket and get it back out.
func ExampleBucket_Put() {
	db, err := bbolt.Open(tempfile(), 0600, &bbolt.Options{
		Timeout: 30 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer os.Remove(db.Path())
	defer db.Close()

	// Create a new `things` bucket.
	bucket := []byte("things")
	things, _ := store.NewBucket(db, bucket)

	// Put key/value into the `things` bucket.
	key, value := []byte("A"), []byte("alpha")
	if err := things.Put(key, value); err != nil {
		fmt.Printf("could not insert item: %v", err)
	}

	// Read value back in a different read-only transaction.
	got, _ := things.Get(key)

	fmt.Printf("The value of %q in `%s` is %q\n", key, bucket, got)

	// Output:
	// The value of "A" in `things` is "alpha"
}

// Show we don't overwrite existing values when using PutNX.
func ExampleBucket_PutNX() {
	db, err := bbolt.Open(tempfile(), 0600, &bbolt.Options{
		Timeout: 30 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer os.Remove(db.Path())
	defer db.Close()

	// Create a new `things` bucket.
	bucket := []byte("things")
	things, _ := store.NewBucket(db, bucket)

	// Put key/value into the `things` bucket.
	key, value := []byte("A"), []byte("alpha")
	if err := things.Put(key, value); err != nil {
		fmt.Printf("could not insert item: %v", err)
	}

	// Read value back in a different read-only transaction.
	got, _ := things.Get(key)

	fmt.Printf("The value of %q in `%s` is %q\n", key, bucket, got)

	// Try putting another value with same key.
	things.PutNX(key, []byte("beta"))

	// Read value back in a different read-only transaction.
	got, _ = things.Get(key)

	fmt.Printf("The value of %q in `%s` is still %q\n", key, bucket, got)

	// Output:
	// The value of "A" in `things` is "alpha"
	// The value of "A" in `things` is still "alpha"
}

// Show we can insert items into a bucket and get them back out.
func ExampleBucket_Insert() {
	db, err := bbolt.Open(tempfile(), 0600, &bbolt.Options{
		Timeout: 30 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer os.Remove(db.Path())
	defer db.Close()

	letters, _ := store.NewBucket(db, []byte("letters"))

	// Setup items to insert in `letters` bucket.
	items := []struct {
		Key, Value []byte
	}{
		{[]byte("A"), []byte("alpha")},
		{[]byte("B"), []byte("beta")},
		{[]byte("C"), []byte("gamma")},
	}

	// Insert items into `letters` bucket.
	if err := letters.Insert(items); err != nil {
		fmt.Println("could not insert items!")
	}

	// Get items back out in separate read-only transaction.
	results, _ := letters.Items()

	for _, item := range results {
		fmt.Printf("%s -> %s\n", item.Key, item.Value)
	}

	// Output:
	// A -> alpha
	// B -> beta
	// C -> gamma
}
