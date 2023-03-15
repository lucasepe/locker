package cmd

import (
	"flag"
	"strings"

	"github.com/lucasepe/locker/cmd/flags"
	"github.com/lucasepe/locker/internal/term"
)

func newCmdList() *cmdList {
	return &cmdList{
		namespace: flags.Namespace{},
		storeRef: flags.Store{
			BaseDir: AppDir(),
		},
	}
}

type cmdList struct {
	namespace flags.Namespace
	storeRef  flags.Store
}

func (*cmdList) Name() string { return "list" }
func (*cmdList) Synopsis() string {
	return "List all namespaces or all keys in a namespace."
}

func (*cmdList) Usage() string {
	return strings.ReplaceAll(`{NAME} list [flags]
  
   List all keys in the namespace 'google' from the store 'accounts':
     {NAME} list -n google -s accounts'

   List all namespaces in the default store:
   {NAME} list`, "{NAME}", appLowerName)
}

func (c *cmdList) SetFlags(fs *flag.FlagSet) {
	fs.Var(&c.namespace, "n", "Namespace.")
	fs.Var(&c.storeRef, "s", "Store name.")
}

func (c *cmdList) Execute(fs *flag.FlagSet) error {
	if len(c.namespace.Bytes()) == 0 {
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

	all, err := db.Namespaces()
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

	all, err := db.Keys(c.namespace.String())
	if err != nil {
		return err
	}
	if len(all) == 0 {
		return nil
	}

	term.PrintColumns(fs.Output(), &all, 6)

	return nil
}
