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

func newCmdDelete() *cmdDelete {
	return &cmdDelete{
		box:      flags.BoxFlag{},
		key:      flags.LabelFlag{},
		storeRef: flags.StoreFlag{},
	}
}

type cmdDelete struct {
	box      flags.BoxFlag
	key      flags.LabelFlag
	storeRef flags.StoreFlag
}

func (*cmdDelete) Name() string { return "delete" }
func (*cmdDelete) Synopsis() string {
	return "Delete one secret from a box or a whole box."
}

func (*cmdDelete) Usage() string {
	return strings.ReplaceAll(`{NAME} delete [flags]
  
   Delete the secret with label 'user' in the box 'Google':
     {NAME} delete -b Google -l user

   Delete the box 'Google':
     {NAME} get -b Google`, "{NAME}", app.Name)
}

func (c *cmdDelete) SetFlags(fs *flag.FlagSet) {
	fs.Var(&c.box, "b", "Box title.")
	fs.Var(&c.storeRef, "n", "Locker name.")
	fs.Var(&c.key, "l", "Secret label.")
}

func (c *cmdDelete) Execute(fs *flag.FlagSet) error {
	if err := c.complete(fs); err != nil {
		return err
	}

	db, err := c.storeRef.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	if len(c.key.Bytes()) > 0 {
		return c.deleteItem(db, fs)
	}

	err = store.DeleteBucket(db, c.box.Bytes())
	if err == nil {
		fmt.Fprintf(fs.Output(), "box successfully deleted (box: %s)\n", c.box.String())
	}

	return nil
}

func (c *cmdDelete) deleteItem(db *bbolt.DB, fs *flag.FlagSet) error {
	bkt, err := store.NewBucket(db, c.box.Bytes())
	if err != nil {
		return err
	}

	err = bkt.Delete(c.key.Bytes())
	if err == nil {
		fmt.Fprintf(fs.Output(), "secret successfully deleted (label: %s, box: %s)\n",
			c.key.String(), c.box.String())
	}

	return nil
}

func (c *cmdDelete) complete(fs *flag.FlagSet) error {
	if len(app.MasterPassword()) == 0 {
		return app.ErrUnsetMasterPassword
	}

	if len(c.box.Bytes()) == 0 {
		return fmt.Errorf("missing: box title")
	}

	return nil
}
