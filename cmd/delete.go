package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lucasepe/locker/cmd/flags"
)

func newCmdDelete() *cmdDelete {
	return &cmdDelete{
		namespace: flags.Namespace{},
		key:       flags.Key{},
		storeRef: flags.Store{
			BaseDir: AppDir(),
		},
	}
}

type cmdDelete struct {
	namespace flags.Namespace
	key       flags.Key
	storeRef  flags.Store
}

func (*cmdDelete) Name() string { return "delete" }
func (*cmdDelete) Synopsis() string {
	return "Delete one or all secrets from a namespace."
}

func (*cmdDelete) Usage() string {
	return strings.ReplaceAll(`{NAME} delete [flags]
  
   Delete the secret with key 'user' in the namespace 'google':
     {NAME} delete -n google -k user

   Delete the 'google' namespace:
     {NAME} delete -n google`, "{NAME}", appLowerName)
}

func (c *cmdDelete) SetFlags(fs *flag.FlagSet) {
	fs.Var(&c.namespace, "n", "Namespace.")
	fs.Var(&c.storeRef, "s", "Store name.")
	fs.Var(&c.key, "k", "Secret key.")
}

func (c *cmdDelete) Execute(fs *flag.FlagSet) error {
	if err := c.complete(fs); err != nil {
		return err
	}

	sto, err := c.storeRef.Connect()
	if err != nil {
		return err
	}
	defer sto.Close()

	if len(c.key.Bytes()) > 0 {
		err := sto.DeleteOne(c.namespace.String(), c.key.String())
		if err == nil {
			fmt.Fprintf(fs.Output(), "secret successfully deleted (key: %s, namespace: %s)\n",
				c.key.String(), c.namespace.String())
		}
		return nil
	}

	if err := sto.DeleteAll(c.namespace.String()); err == nil {
		fmt.Fprintf(fs.Output(), "namespace '%s' successfully deleted\n", c.namespace.String())
	}

	return nil
}

func (c *cmdDelete) complete(fs *flag.FlagSet) error {
	if len(c.namespace.Bytes()) == 0 {
		return fmt.Errorf("missing namespace")
	}

	return nil
}
