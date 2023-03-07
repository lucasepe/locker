package cmd

import (
	"flag"
	"strings"

	"github.com/lucasepe/locker/cmd/app"
	"github.com/lucasepe/locker/cmd/flags"
	"github.com/lucasepe/locker/internal/store"
	"github.com/lucasepe/locker/internal/term"
)

func newCmdList() *cmdList {
	return &cmdList{
		box:      flags.BoxFlag{},
		storeRef: flags.StoreFlag{},
	}
}

type cmdList struct {
	box      flags.BoxFlag
	storeRef flags.StoreFlag
}

func (*cmdList) Name() string { return "list" }
func (*cmdList) Synopsis() string {
	return "List all boxes or all items in a box."
}

func (*cmdList) Usage() string {
	return strings.ReplaceAll(`{NAME} list [flags]
  
   List all items in the box 'Google' from the locker 'accounts':
     {NAME} list -b Google -n accounts'

   List all boxes in the default locker:
   {NAME} list`, "{NAME}", app.Name)
}

func (c *cmdList) SetFlags(fs *flag.FlagSet) {
	fs.Var(&c.box, "b", "Box title.")
	fs.Var(&c.storeRef, "n", "Locker name.")
}

func (c *cmdList) Execute(fs *flag.FlagSet) error {
	if len(c.box.Bytes()) == 0 {
		return c.printBuckets(fs)
	}

	return c.printKeysInBucket(fs)
}

func (c *cmdList) printBuckets(fs *flag.FlagSet) error {
	db, err := c.storeRef.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	all, err := store.Buckets(db)
	if err != nil {
		return err
	}
	if len(all) > 0 {
		term.PrintColumns(fs.Output(), &all, 6)
	}

	return nil
}

func (c *cmdList) printKeysInBucket(fs *flag.FlagSet) error {
	db, err := c.storeRef.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	bkt, err := store.NewBucket(db, c.box.Bytes())
	if err != nil {
		return err
	}

	all, err := bkt.Keys()
	if err != nil {
		return err
	}
	if len(all) == 0 {
		return nil
	}

	term.PrintColumns(fs.Output(), &all, 6)

	return nil

}
