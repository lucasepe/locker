package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lucasepe/locker/cmd/app"
	"github.com/lucasepe/locker/cmd/flags"
	"github.com/lucasepe/locker/internal/store"
	"go.etcd.io/bbolt"
)

func newCmdGet() *cmdGet {
	return &cmdGet{
		box:      flags.BoxFlag{},
		key:      flags.LabelFlag{},
		storeRef: flags.StoreFlag{},
	}
}

type cmdGet struct {
	box      flags.BoxFlag
	key      flags.LabelFlag
	storeRef flags.StoreFlag
}

func (*cmdGet) Name() string { return "get" }
func (*cmdGet) Synopsis() string {
	return "Get one or all secrets from a box."
}

func (*cmdGet) Usage() string {
	return strings.ReplaceAll(`{NAME} get [flags]
  
   Get the secret with label 'user' from the 'Google' box:
     {NAME} get -b Google -l user

   Get all the secrets from the 'Google' box:
     {NAME} get -b Google`, "{NAME}", app.Name)
}

func (c *cmdGet) SetFlags(fs *flag.FlagSet) {
	fs.Var(&c.box, "b", "Box title.")
	fs.Var(&c.storeRef, "n", "Locker name.")
	fs.Var(&c.key, "l", "Secret label.")
}

func (c *cmdGet) Execute(fs *flag.FlagSet) error {
	if err := c.complete(fs); err != nil {
		return err
	}

	db, err := c.storeRef.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	if len(c.key.Bytes()) > 0 {
		return c.extractItem(db, fs)
	}

	return c.extractAll(db, fs)
}

func (c *cmdGet) extractAll(db *bbolt.DB, fs *flag.FlagSet) error {
	bkt, err := store.NewBucket(db, c.box.Bytes())
	if err != nil {
		return err
	}

	all := map[string]string{}
	fn := func(k, v []byte) error {
		dec, err := app.Decrypt(v)
		if err != nil {
			return err
		}

		all[string(k)] = string(dec)
		return nil
	}

	if err := bkt.ForEach(fn); err != nil {
		return err
	}

	for k, v := range all {
		fmt.Fprintf(fs.Output(), "%s: %s\n", k, v)
	}

	return nil
}

func (c *cmdGet) extractItem(db *bbolt.DB, fs *flag.FlagSet) error {
	bkt, err := store.NewBucket(db, c.box.Bytes())
	if err != nil {
		return err
	}

	dat, err := bkt.Get(c.key.Bytes())
	if err != nil {
		return err
	}
	if dat == nil {
		return nil
	}

	val, err := app.Decrypt(dat)
	if err != nil {
		return err
	}

	fmt.Fprintf(fs.Output(), "%s\n", val)

	return nil
}

func (c *cmdGet) complete(fs *flag.FlagSet) error {
	if len(app.MasterPassword()) == 0 {
		return app.ErrUnsetMasterPassword
	}

	if len(c.box.Bytes()) == 0 {
		return fmt.Errorf("missing: box title")
	}

	return nil
}
