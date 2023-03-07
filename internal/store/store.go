package store

import (
	"go.etcd.io/bbolt"
)

// Bucket represents a collection of key/value pairs inside the database.
type Bucket struct {
	db   *bbolt.DB
	Name []byte
}

// An Item holds a key/value pair.
type Item struct {
	Key   []byte
	Value []byte
}

// NewBucket creates/opens a named bucket.
func NewBucket(db *bbolt.DB, name []byte) (*Bucket, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Bucket{db: db, Name: name}, nil
}

// DeleteBucket removes the named bucket.
func DeleteBucket(db *bbolt.DB, name []byte) error {
	return db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket(name)
	})
}

// Buckets return all buckets names.
func Buckets(db *bbolt.DB) (names []string, err error) {
	err = db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(bn []byte, _ *bbolt.Bucket) error {
			el := make([]byte, len(bn))
			copy(el, bn)
			names = append(names, string(el))
			return nil
		})
	})

	return names, err
}

// Put inserts value `v` with key `k`.
func (bk *Bucket) Put(k, v []byte) error {
	return bk.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bk.Name).Put(k, v)
	})
}

// PutNX (put-if-not-exists) inserts value `v` with key `k`
// if key doesn't exist.
func (bk *Bucket) PutNX(k, v []byte) error {
	old, err := bk.Get(k)
	if old != nil || err != nil {
		return err
	}
	return bk.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bk.Name).Put(k, v)
	})
}

// Insert iterates over a slice of k/v pairs, putting each item in
// the bucket as part of a single transaction.  For large insertions,
// be sure to pre-sort your items (by Key in byte-sorted order), which
// will result in much more efficient insertion times and storage costs.
func (bk *Bucket) Insert(items []struct{ Key, Value []byte }) error {
	return bk.db.Update(func(tx *bbolt.Tx) error {
		for _, item := range items {
			tx.Bucket(bk.Name).Put(item.Key, item.Value)
		}
		return nil
	})
}

// InsertNX (insert-if-not-exists) iterates over a slice of k/v pairs,
// putting each item in the bucket as part of a single transaction.
// Unlike Insert, however, InsertNX will not update the value for an
// existing key.
func (bk *Bucket) InsertNX(items []struct{ Key, Value []byte }) error {
	return bk.db.Update(func(tx *bbolt.Tx) error {
		for _, item := range items {
			v, _ := bk.Get(item.Key)
			if v == nil {
				tx.Bucket(bk.Name).Put(item.Key, item.Value)
			}
		}
		return nil
	})
}

// Delete removes key `k`.
func (bk *Bucket) Delete(k []byte) error {
	return bk.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bk.Name).Delete(k)
	})
}

// Get retrieves the value for key `k`.
func (bk *Bucket) Get(k []byte) (value []byte, err error) {
	err = bk.db.View(func(tx *bbolt.Tx) error {
		v := tx.Bucket(bk.Name).Get(k)
		if v != nil {
			value = make([]byte, len(v))
			copy(value, v)
		}
		return nil
	})
	return value, err
}

func (bk *Bucket) ForEach(fn func(k, v []byte) error) error {
	return bk.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket(bk.Name).Cursor()
		for key, val := c.First(); key != nil; key, val = c.Next() {
			k := make([]byte, len(key))
			copy(k, key)

			v := make([]byte, len(val))
			copy(v, val)

			if err := fn(k, v); err != nil {
				return err
			}
		}

		return nil
	})
}

// Items returns a slice of key/value pairs.  Each k/v pair in the slice
// is of type Item (`struct{ Key, Value []byte }`).
func (bk *Bucket) Items() (items []Item, err error) {
	return items, bk.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket(bk.Name).Cursor()
		var key, value []byte
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				key = make([]byte, len(k))
				copy(key, k)
				value = make([]byte, len(v))
				copy(value, v)
				items = append(items, Item{key, value})
			}
		}
		return nil
	})
}

// Keys returns all keys in the bucket.
func (bk *Bucket) Keys() (items []string, err error) {
	return items, bk.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket(bk.Name).Cursor()

		var key []byte
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			key = make([]byte, len(k))
			copy(key, k)

			items = append(items, string(key))
		}
		return nil
	})
}
