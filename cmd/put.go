package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lucasepe/locker/cmd/app"
	"github.com/lucasepe/locker/cmd/flags"
	"github.com/lucasepe/locker/internal/store"
)

func newCmdPut() *cmdPut {
	return &cmdPut{
		box:      flags.BoxFlag{},
		key:      flags.LabelFlag{},
		storeRef: flags.StoreFlag{},
	}
}

type cmdPut struct {
	box      flags.BoxFlag
	key      flags.LabelFlag
	storeRef flags.StoreFlag
}

func (*cmdPut) Name() string { return "put" }
func (*cmdPut) Synopsis() string {
	return "Put an item into a box."
}

func (*cmdPut) Usage() string {
	return strings.ReplaceAll(`{NAME} put [flags]
  
   Put the item with label 'user' and value 'my@gmail.com' into the 'Google' box:
     {NAME} put -b Google -l user my@gmail.com'

   Put the content of the 'doc.txt' file as label 'myDoc' into the box 'Docs':
     cat doc.txt | {NAME} put -b Docs -l myDoc

   Put an item whose content is another command output (using pipes):
     pwgen 14 1 | {NAME} put -b Instagram -l pass`, "{NAME}", app.Name)
}

func (c *cmdPut) SetFlags(fs *flag.FlagSet) {
	fs.Var(&c.box, "b", "Box title.")
	fs.Var(&c.storeRef, "n", "Locker name.")
	fs.Var(&c.key, "l", "Item label.")
}

func (c *cmdPut) Execute(fs *flag.FlagSet) error {
	if err := c.complete(fs); err != nil {
		return err
	}

	val := grabContent(fs)
	if len(val) == 0 {
		return nil
	}

	val, err := app.Encrypt(val)
	if err != nil {
		return err
	}

	db, err := c.storeRef.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	bkt, err := store.NewBucket(db, c.box.Bytes())
	if err != nil {
		return err
	}

	err = bkt.Put(c.key.Bytes(), val)
	if err == nil {
		fmt.Fprintf(fs.Output(), "item successfully stored (label:%s, box: %s, locker: %s)\n",
			c.key.String(), c.box.String(), c.storeRef.Name())
	}

	return err
}

func (c *cmdPut) complete(fs *flag.FlagSet) error {
	if len(app.MasterPassword()) == 0 {
		return app.ErrUnsetMasterPassword
	}

	if len(c.box.Bytes()) == 0 {
		return fmt.Errorf("missing: box title")
	}

	if len(c.key.Bytes()) == 0 {
		return fmt.Errorf("missing: memo key")
	}

	return nil
}
