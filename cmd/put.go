package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lucasepe/locker/cmd/app"
	"github.com/lucasepe/locker/cmd/flags"
)

func newCmdPut() *cmdPut {
	return &cmdPut{
		namespace: flags.NamespaceFlag{},
		key:       flags.KeyFlag{},
		storeRef:  flags.StoreFlag{},
	}
}

type cmdPut struct {
	namespace flags.NamespaceFlag
	key       flags.KeyFlag
	storeRef  flags.StoreFlag
}

func (*cmdPut) Name() string { return "put" }
func (*cmdPut) Synopsis() string {
	return "Put a secret into a namespace."
}

func (*cmdPut) Usage() string {
	return strings.ReplaceAll(`{NAME} put [flags]
  
   Put the secret with key 'user' and value 'my@gmail.com' into the 'google' namespace:
     {NAME} put -n google -k user my@gmail.com'

   Put the content of the 'doc.txt' file with key 'my_doc' into the namespace 'docs':
     cat doc.txt | {NAME} put -n docs -k my_doc

   Put a secret whose content is another command output (using pipes):
     pwgen 14 1 | {NAME} put -n Instagram -k password`, "{NAME}", app.Name)
}

func (c *cmdPut) SetFlags(fs *flag.FlagSet) {
	fs.Var(&c.namespace, "n", "Namespace.")
	fs.Var(&c.storeRef, "s", "Store name.")
	fs.Var(&c.key, "k", "Secret key.")
}

func (c *cmdPut) Execute(fs *flag.FlagSet) error {
	if err := c.complete(fs); err != nil {
		return err
	}

	val := grabContent(fs)
	if len(val) == 0 {
		return nil
	}

	sto, err := c.storeRef.Connect()
	if err != nil {
		return err
	}
	defer sto.Close()

	err = sto.PutOne(c.namespace.String(), c.key.String(), string(val))
	if err == nil {
		fmt.Fprintf(fs.Output(), "secret successfully stored (key:%s, namespace: %s, store: %s)\n",
			c.key.String(), c.namespace.String(), c.storeRef.Name())
	}

	return err
}

func (c *cmdPut) complete(fs *flag.FlagSet) error {
	if len(c.namespace.Bytes()) == 0 {
		return fmt.Errorf("missing namespace")
	}

	if len(c.key.Bytes()) == 0 {
		return fmt.Errorf("missing key")
	}

	return nil
}
